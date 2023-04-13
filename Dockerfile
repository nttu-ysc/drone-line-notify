FROM golang:1.20.3-alpine

RUN mkdir -p /app

WORKDIR /app

COPY . .

RUN go build -o main .

CMD [ "/app/main" ]
