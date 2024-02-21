# go build env
FROM golang:1.20.3-alpine AS build-env

RUN mkdir -p /app

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o main .

# run env
FROM alpine

ENV TZ="Asia/Taipei"

RUN apk update && apk add -U tzdata

WORKDIR /app

COPY --from=build-env /app /app

ENTRYPOINT [ "/app/main" ]