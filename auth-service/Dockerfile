FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd/main.go

#----
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/auth-service .
COPY .env .

EXPOSE 8080
CMD [ "./auth-service" ]