FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd/main.go

# ----
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/auth-service .
COPY migrations /app/migrations
COPY .env .

EXPOSE 8080 50051
CMD ["./auth-service"]
