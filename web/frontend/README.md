# cerbero-frontend

> 📦 Frontend of cerbero-web.

## Table of contents

<!--toc:start-->
- [cerbero-frontend](#cerbero-frontend)
  - [Table of contents](#table-of-contents)
  - [Developers documentation](#developers-documentation)
    - [Quick start](#quick-start)
      - [Install the required `node_modules` for development](#install-the-required-nodemodules-for-development)
      - [Start the development server](#start-the-development-server)
    - [Build](#build)
      - [Install the required `node_modules` for building](#install-the-required-nodemodules-for-building)
      - [Compile the source into javascript](#compile-the-source-into-javascript)
      - [Finally preview the compiled source with](#finally-preview-the-compiled-source-with)
<!--toc:end-->

## Developers documentation

### Quick start

> The following commands **must** be executed inside the `web/frontend/` directory.

#### Install the required `node_modules` for development

```sh
npm i
```

#### Start the development server

```sh
npm run dev
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

#### Finally preview the compiled source with

```
npm run preview
```

