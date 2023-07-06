# Build stage
FROM golang:1.20-alpine3.17 as builder

WORKDIR /app/go-bank

COPY . .

RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.17 

WORKDIR /app/go-bank

COPY --from=builder /app/go-bank/main .
COPY --from=builder /app/go-bank/migrate ./migrate
COPY app.env .
COPY db/migration ./migration
COPY start.sh .

EXPOSE 5000

CMD ["/app/go-bank/main"]
ENTRYPOINT ["/app/go-bank/start.sh"]

