# Build stage
FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server .

# Runtime stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/server .
ENV VERSION="v1"
EXPOSE 8080
CMD ["./server"]