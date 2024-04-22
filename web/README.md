# cerbero-web

> 🌐 Web-based control panel for cerbero.

## Table of contents

<!--toc:start-->
- [cerbero-web](#cerbero-web)
  - [Table of contents](#table-of-contents)
  - [Users documentation](#users-documentation)
    - [Deployment with docker compose](#deployment-with-docker-compose)
      - [Environment variables](#environment-variables)
      - [Start the containers](#start-the-containers)
  - [Developers documentation](#developers-documentation)
    - [Develop the frontend](#develop-the-frontend)
    - [Develop the backend](#develop-the-backend)
<!--toc:end-->

## Users documentation

### Deployment with docker compose

#### Environment variables

Create a `.env` file from the `.template.env`:

```sh
cp .template.env .env
```

Fill the `.env` file with the following values:

```
API_PORT="80"
REDIS_URL="redis://cerbero-redis-stack:6379"
SOCKET_PORT="6969"
```

#### Start the containers

```sh
docker compose up -d
```

## Developers documentation

### Develop the frontend

You can find more information about the frontend [here](/web/frontend/README.md#developers-documentation).

### Develop the backend

You can find more information about the backend [here](/web/backend/README.md#developers-documentation).

