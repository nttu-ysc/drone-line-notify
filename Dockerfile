FROM golang:1.20.3-alpine

RUN mkdir -p /app

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o main .

ENTRYPOINT [ "/app/main"]
