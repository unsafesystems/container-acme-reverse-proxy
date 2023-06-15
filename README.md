# ACME Reverse Proxy

This repository provides a GoLang / Docker based, ACME-enabled reverse proxy. The goal was to provide a very
simple, easy to use, reverse proxy that can be used to front web applications.

## Configuration

The reverse proxy is configured via environment variables. The following variables are supported:

| Variable       | Description                                                | Example                 |
|----------------|------------------------------------------------------------|-------------------------|
| `PROXY_HOST`   | The host the proxy is listening on. Used with letsencrypt. | `example.com`           |
| `PROXY_EMAIL`  | The email address to register with letsencrypt.            | `mail@example.com`      |
| `PROXY_TARGET` | The target URL to proxy requests to.                       | `http://localhost:8080` |


## Usage

The following example shows how to run the reverse proxy:

```bash
docker run --rm -p 8080:8080 -p 8443:8443 \
  --env PROXY_HOST=76f1-45-14-97-5.ngrok-free.app \
  --env PROXY_EMAIL=ferdinand@vonhagen.me \
  --env PROXY_TARGET=http://localhost:3000 \
  -v "$(pwd)/certs":/home/autocert \
  acme-proxy
```

The proxy will listen on port 8080 for HTTP requests and on port 8443 for HTTPS requests.

Mounting a (persistent) volume to `/home/autocert` is recommended to persist the certificates. Let's Encrypt restricts
the number of certificates that can be issued per week. Persisting the certificates allows the proxy to reuse them
across restarts.