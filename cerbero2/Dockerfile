FROM golang:1.21

WORKDIR /app
COPY . .
RUN go mod tidy

CMD go build -o /app/bin/firewall2
