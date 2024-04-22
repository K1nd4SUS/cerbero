# cerbero-backend

> 📦 Backend of cerbero-web.

## Table of contents

<!--toc:start-->
- [cerbero-backend](#cerbero-backend)
  - [Table of contents](#table-of-contents)
  - [Developers documentation](#developers-documentation)
    - [Quick start](#quick-start)
      - [Install the required `node_modules` for development](#install-the-required-nodemodules-for-development)
      - [Create a `.env` file](#create-a-env-file)
      - [Start a redis-stack instance](#start-a-redis-stack-instance)
      - [Start the development server](#start-the-development-server)
    - [Build](#build)
      - [Install the required `node_modules` for building](#install-the-required-nodemodules-for-building)
      - [Compile the source into javascript](#compile-the-source-into-javascript)
      - [Finally start the compiled source with](#finally-start-the-compiled-source-with)
    - [Environment variables](#environment-variables)
      - [`API_PORT`](#apiport)
      - [`REDIS_URL`](#redisurl)
      - [`SOCKET_PORT`](#socketport)
<!--toc:end-->

## Developers documentation

### Quick start

> The following commands **must** be executed inside the `web/backend/` directory.

*Make sure to be on node version `18.17.1`, you can quickly swap between node versions with [`nvm`](https://github.com/nvm-sh/nvm) (node version manager).*

#### Install the required `node_modules` for development

```sh
npm i
```

#### Create a `.env` file

> The `.env` file **must** be located inside the `web/backend/` directory (`web/backend/.env`).

Fill the `.env` with the following values:

```sh
API_PORT="6666"
REDIS_URL="redis://localhost:6379"
SOCKET_PORT="6969"
```

**Currently the `API_PORT` value MUST be `6666` because in the development phase `/api` requests are proxied only to that port (it is hardcoded in `vite.config.ts`)**

#### Start a redis-stack instance

One of the easiest ways to bootstrap a redis-stack instance is running it with docker:

```sh
docker run -d -p "127.0.0.1:6379:6379" --name cerbero-redis-stack-dev redis/redis-stack:latest
```

#### Start the development server

```sh
npm run dev
```

If you have done everything correctly you should be able to see something like this in your terminal:

```sh
2012-12-12T12:00:00.000Z [INFO] API listening on port 6666
2012-12-12T12:00:00.000Z [INFO] Socket server listening on port 6969
2012-12-12T12:00:00.000Z [INFO] Connected to db redis://localhost:6379
```

### Build

> The build process consists in compiling the typescript code into javascript, the compiled source will be stored into the `dist/` directory.

*Most of the times you won't need to build anything manually, we use automated CI pipelines to handle that for us.*

#### Install the required `node_modules` for building

```sh
npm i
```

#### Compile the source into javascript

```sh
npm run build
```

#### Finally start the compiled source with

```
npm run start
```

### Environment variables

#### `API_PORT`

This variable is **mandatory** and specifies the port where the express api will start listening for incoming requests.

#### `REDIS_URL`

This variable is **mandatory** and specifies the redis connection string that the api will use to connect to the database.

#### `SOCKET_PORT`

This variable is **mandatory** and specifies the port where the TCP socket server will start listening for connections.

