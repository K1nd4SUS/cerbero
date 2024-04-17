FROM node:18.17.1-slim

WORKDIR /app

COPY package* .
COPY dist dist

RUN npm i --omit=dev

CMD ["npm", "run", "start"]

