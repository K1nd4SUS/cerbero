FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod go.sum config.json ./
RUN go mod download

COPY *.go ./

RUN go build -o /cerbero

EXPOSE 8082

RUN apk add --no-cache sudo

CMD ["sudo", "/cerbero" ]