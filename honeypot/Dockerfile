FROM golang:1.25.5 AS builder
WORKDIR /app
COPY honeypot/go.mod ./
RUN go mod download
COPY honeypot/ ./
RUN CGO_ENABLED=0 go build -o honeypot .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/honeypot .
CMD ["./honeypot"]