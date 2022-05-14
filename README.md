# lego RWTH DNS provider

__A [lego](https://github.com/go-acme/lego) provider for [RWTH DNS API](https://noc-portal.rz.rwth-aachen.de/dns-admin/en/api_tokens)__

[Traefik](https://github.com/traefik/traefik) uses lego as their ACME client, so this can be directly used in Traefik.

It is written in Go to allow easy integration with the Traefik docker image which does not include cURL or python.
