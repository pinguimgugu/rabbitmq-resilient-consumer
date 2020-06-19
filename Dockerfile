FROM golang:1.13

WORKDIR /go/src/rabbitmq-consumer

COPY . .

RUN go mod init