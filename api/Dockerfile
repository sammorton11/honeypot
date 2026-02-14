FROM golang:1.25.5 AS builder
WORKDIR /app
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ ./
RUN CGO_ENABLED=0 go build -o api ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api .
CMD ["./api"]