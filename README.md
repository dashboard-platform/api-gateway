# API Gateway

![coverage](https://img.shields.io/badge/coverage-0%25-green) ![go-report](https://goreportcard.com/badge/github.com/dashboard-platform/api-gateway) [![Go Reference](https://pkg.go.dev/badge/github.com/dashboard-platform/api-gateway.svg)](https://pkg.go.dev/github.com/dashboard-platform/api-gateway)


This service acts as the entrypoint to all client traffic and handles request forwarding to individual microservices such as `auth-service` and `dashboard-service`. It also manages authentication, routing, and headers.

```bash
git clone https://github.com/dashboard-platform/api-gateway.git
cd api-gateway
```

## Run Locally

### Option 1: Run with Go (requires PostgreSQL running locally)

1. Create a `.env` file.
2. Run:

```bash
go run cmd/main.go
```

`.env` file example: (for details, see below)
```env
PORT=:8080
AUTH_SERVICE_URL=http://auth-service:8080
DASHBOARD_SERVICE_URL=http://dashboard-service:8080
JWT_SECRET=supersecretkey
COOKIE_SECURE=false

```

### Option 2: Run with Docker

```bash
docker build -t api-gateway .
docker run -p 8080:8080 --env-file .env api-gateway
```

This will start the `api-gateway` on port `8080`.

Access healthcheck:
```bash
curl http://localhost:8080/healthcheck
```

## Run Tests

To run the whole test suite, use:
```bash
cd api-gateway
go tests -v ./...
```

## Environment Variables

| Variable     | Description                             |
|--------------|-----------------------------------------|
| `PORT`       | Port on which the service runs (`:8080`)          |
| `AUTH_SERVICE`     | Address where `auth-service` is running  |
| `DASHBOARD_SERVICE` | Address where `dashboard-service` is running |
| `JWT_SECRET` | Secret used for signing JWTs (`secret`)        |
| `COOKIE_SECURE`        | Use secured cookies or not |

## Features

- Centralized routing for all internal APIs
- JWT validation middleware (from cookie or bearer token)
- Forwarding to internal microservices:
  - `/auth/*` → `auth-service`
  - `/dashboard/*` → `dashboard-service`
- Cookie handling and header normalization
- Built-in support for CORS and secure HTTP headers

## Endpoints

| Method | Path         | Auth Required | Description                       |
|--------|--------------|----------------|-----------------------------------|
| GET    | `/healthcheck` | ❌             | Basic service |  