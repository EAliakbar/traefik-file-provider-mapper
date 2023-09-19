<h1 align="center">Welcome to traefik-file-provider-mapper 👋</h1>

> A traefik provider mapping server to edit an existing file provider, currently focused on integrating Coolify V4 into an existing traefik using http provider.
### 🏠 [Homepage](https://github.com/EAliakbar/traefik-file-provider-mapper#readme)
### ✨ [Blog Post (Coming Soon)]()

> Coolify V4 have switched from traefik http provider to file provider  
> There is a [traefik-http-provider-mapper](https://github.com/beaussan/traefik-http-provider-mapper) written by @beaussan to integrate Coolify V3 into an existing traefik  
> This Project aims to do the same for Coolify V4

## Docker compose install

```yaml
version: '3'

services:
  traefik-file-provider-mapper:
    image: ghcr.io/ealiakbar/traefik-file-provider-mapper
    restart: unless-stopped
    volumes:
      - "/data/coolify/proxy/dynamic/:/input"
    environment:
      # All the following are optional

      # Replace the web entrypoint with what is provided (eg http here)
      - TRAEFIK_MAPPER_NEW_HTTP_ENTRYPOINT=http
      # Replace the websecure entrypoint with what is provided (eg https here)
      - TRAEFIK_MAPPER_NEW_HTTPS_ENTRYPOINT=https
      # Replace the certificate store (letsencrypt) with what is provided (eg cert here)
      - TRAEFIK_MAPPER_NEW_CERT_RESOLVER=cert

      # Adds this middleware to all routes, ex middlewareName@file for a file base middleware, middlewareName@docker for a docker base middleware
      - TRAEFIK_MAPPER_ADD_MIDDLEWARE_NAME=middlewareName@source

      # Ignore routers for the middlewares. This can be usefull to filter out some domains
      - TRAEFIK_MAPPER_IGNORE_MIDDLEWARE_SITES=some.domain.io;some.other.domain.io

      # Remove coolify itself from the list of services
      - TRAEFIK_MAPPER_FILTER_COOLIFY=true
      # Remove www redirect middleware from the list of middlewares
      - TRAEFIK_MAPPER_FILTER_WWW_MIDDLEWARE=true
      # Remove all http routes from the list of routes
      - TRAEFIK_MAPPER_FILTER_HTTP_ROUTERS=true
      # Remove all www routes from the list of routes
      - TRAEFIK_MAPPER_FILTER_WWW_ROUTERS=true

    networks:
      # Hooks itself into the Traefik network
      - coolify
networks:
  coolify:
    external: true

```

And then, in your Traefik config, add this url as a [http provider](https://doc.traefik.io/traefik/providers/http/#provider-configuration)

> http://traefik-file-provider-mapper:8080/

And then, you should have every Coolify provided routes in your Traefik instance ! :tada:

## Limitations
- `TRAEFIK_MAPPER_FILTER_WWW_MIDDLEWARE` is not yet supported
- `TRAEFIK_MAPPER_ADD_MIDDLEWARE_NAME` is not yet supported
- `TRAEFIK_MAPPER_IGNORE_MIDDLEWARE_SITES` is not yet supported