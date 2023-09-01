FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o /test_service

EXPOSE 8083

CMD [ "/test_service" ]