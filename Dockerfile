FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api-gateway ./cmd/main.go

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/api-gateway .

EXPOSE 8080

# Run the binary
CMD ["./api-gateway"]
