FROM golang:alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

EXPOSE 9090

COPY . .

RUN go mod tidy

RUN go build -o binary

ENTRYPOINT ["/app/binary"]