# cerbero-web

<!--toc:start-->
- [cerbero-web](#cerbero-web)
  - [Deploy](#deploy)
    - [Environment variables](#environment-variables)
    - [Start the containers](#start-the-containers)
<!--toc:end-->

> Web-based control panel for cerbero.

## Deploy

### Environment variables

Create a `.env` file from the `.template.env`:

```sh
cp .template.env .env
```

Fill the `.env` file with the desired values:

Example:

```
API_PORT="8080"
REDIS_URL="redis://cerbero-redis-stack:6379"
SOCKET_PORT="6969"
```

### Start the containers

```sh
docker compose up -d
```

