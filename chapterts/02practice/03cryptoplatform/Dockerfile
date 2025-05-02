# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server main.go

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/server /app/server

EXPOSE 8080

CMD ["/app/server"]
