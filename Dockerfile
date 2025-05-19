# Build stage
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o client-app ./main.go


FROM alpine:latest
WORKDIR /
COPY --from=builder /app/client-app /client-app
RUN chmod +x /client-app
ENTRYPOINT ["/client-app"]