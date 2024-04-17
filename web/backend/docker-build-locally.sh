#!/bin/sh

# Build backend
npm i
npm run build

# Build docker image
docker build . --tag cerbero-backend:local

