# Dockerfile for clean deployment
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o booking-server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/booking-server .
COPY --from=builder /app/index.html .

# Set environment variable to indicate Docker environment
ENV DOCKER_ENV=true
ENV GIN_MODE=release

EXPOSE 8080
CMD ["./booking-server"]