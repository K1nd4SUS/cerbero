#!/bin/sh

# Build frontend
npm i
npm run build

# Build docker image
docker build . --tag cerbero-frontend:local

