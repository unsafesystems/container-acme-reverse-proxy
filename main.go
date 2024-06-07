package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

const (
	ReadTimeout  = 10 * time.Second
	WriteTimeout = 10 * time.Second
	IdleTimeout  = 120 * time.Second
)

// newServer creates a new HTTP server with the given port and handler.
func newServer(port int, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}
}

func main() {
	target := os.Getenv("PROXY_TARGET")
	if target == "" {
		log.Fatal("PROXY_TARGET environment variable not set")
	}

	rpURL, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}

	host := os.Getenv("PROXY_HOST")
	if host == "" {
		log.Fatal("PROXY_HOST environment variable not set")
	}
	hosts := strings.Split(host, ",")

	email := os.Getenv("PROXY_EMAIL")
	if email == "" {
		log.Fatal("PROXY_EMAIL environment variable not set")
	}

	// Setup autocert to retrieve certificates from Let's Encrypt.
	manager := autocert.Manager{
		Cache:      autocert.DirCache("autocert"),
		Prompt:     autocert.AcceptTOS,
		Email:      email,
		HostPolicy: autocert.HostWhitelist(hosts...),
	}

	// Have autocert listen on port 80 to handle HTTP-01 challenges. This will also take care of redirecting
	// HTTP requests to HTTPS.
	server := newServer(8080, manager.HTTPHandler(nil))
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Proxy to the downstream server.
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetXForwarded()
			r.SetURL(rpURL)
		},
	}

	// Create a TLS-enabled server that listens on port 443 and uses the autocert.Manager to retrieve certificates.
	proxyServer := newServer(8443, proxy)
	proxyServer.TLSConfig = manager.TLSConfig()

	log.Fatal(proxyServer.ListenAndServeTLS("", ""))
}
