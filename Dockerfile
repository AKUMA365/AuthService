FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY internal/migrations ./migrations

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/authservic/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/auth-service .
COPY --from=builder /app/config ./config

COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./auth-service"]