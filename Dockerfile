#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/bin/app
COPY . .
RUN go get github.com/go-sql-driver/mysql
CMD go run .
EXPOSE 3000