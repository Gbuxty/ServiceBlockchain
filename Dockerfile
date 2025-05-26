FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /tickers ./cmd/main.go

FROM alpine:latest
RUN mkdir -p /app/logs && chmod 777 /app/logs
WORKDIR /app
COPY --from=builder /tickers .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/dev.yml ./

EXPOSE 8080
CMD ["./tickers"]