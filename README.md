# lego RWTH DNS provider

__A [lego](https://github.com/go-acme/lego) provider for [RWTH DNS API](https://noc-portal.rz.rwth-aachen.de/dns-admin/en/api_tokens)__

[Traefik](https://github.com/traefik/traefik) uses lego as their ACME client, so this can be directly used in Traefik.

It is written in Go to allow easy integration with the Traefik docker image which does not include cURL or python.

## Usage

In addition to the required environment variables listed [here](https://go-acme.github.io/lego/dns/exec/), you also have to provide a variable `RWTH_DNS_API_TOKEN`.

Use it in your `docker-compose.yml` with Traefik like this:

```yaml
services:
  traefik:
    image: traefik
    command:
      ...
      - "--certificatesresolvers.le-dns.acme.dnschallenge=true"
      - "--certificatesresolvers.le-dns.acme.dnschallenge.provider=exec"
      #- "--certificatesresolvers.le-dns.acme.caserver=https://acme-staging-v02.api.letsencrypt.org/directory"
      - "--certificatesresolvers.le-dns.acme.email=postmaster@example.com"
      - "--certificatesresolvers.le-dns.acme.storage=/letsencrypt/acme.json"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /srv/traefik/letsencrypt:/letsencrypt
      - /srv/traefik/providers:/providers
    environment:
      - EXEC_PATH=/providers/lego-rwth-dns-provider
      - EXEC_POLLING_INTERVAL=60 # every 1 minute
      - EXEC_PROPAGATION_TIMEOUT=1200 # 20 minutes
      - RWTH_DNS_API_TOKEN=my-token
```
