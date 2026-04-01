# Gox-Boilerplate
A production-ready, scalable Go boilerplate designed for high-concurrency environments. This project implements a Clean Architecture pattern with a focus on decoupling, strict type-safety via SQLC, and asynchronous processing using RabbitMQ.

## 🏗 Architecture Overview
This project follows Uncle Bob's Clean Architecture to ensure the business logic remains independent of frameworks, UI, or databases.

- cmd/: Entry points for the API Server and Background Worker.
- internal/domain/: Core Business Entities and Repository/Usecase Interfaces (The "Heart" of the app).
- internal/usecase/: Pure Business Logic orchestration.
- internal/repository/: Data persistence implementations (SQLC, Redis, etc.).
- internal/delivery/: Communication layers (HTTP Gin Handlers & RabbitMQ Consumers).
- internal/bootstrap/: Dependency Injection management using Uber-fx.

## 🛠 Tech Stack
Web Framework: Gin-gonic

- Database: PostgreSQL + SQLC (Type-safe SQL compiler)
- Migrations: Golang-migrate (Versioned SQL migrations)
- Caching: Redis (Used for Cache-aside pattern & JWT Blacklisting)
- Message Broker: RabbitMQ (Asynchronous worker for distributed tasks)
- Dependency Injection: Uber-fx
- Logging: Uber-zap (Structured, high-performance logging)
- Config Management: Viper (.yaml, .env support)
- API Documentation: Swaggo (Swagger UI)

## 🚀 Getting Started
1. Prerequisites

- Go 1.25+
- Docker & Docker Compose
- SQLC (go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest)
- Swag (go install github.com/swaggo/swag/cmd/swag@latest)

2. Infrastructure Setup
Spin up the infrastructure services (Postgres, Redis, RabbitMQ) + observability stack (Prometheus, Loki, Grafana):

```bash
make docker-up
```

3. Initialize Database & Code Generation

```bash
# Run SQL migrations
make migrate-up

# Generate Type-safe Go code from SQL (Crucial)
make sqlc

# Generate Swagger documentation
make swag

# Tidy Go modules
make tidy
```

4. Running the Application
Open two separate terminals for the API and the Worker:

```bash
# Start API Server (Default Port: 8080)
make run-api

# Start Background Worker (RabbitMQ Consumer)
make run-worker
```

## 📖 API Documentation
Once the API server is running, you can explore and test the endpoints via Swagger UI:
http://localhost:8080/swagger/index.html

## 📈 Observability (Prometheus, Loki, Grafana)

1. Prometheus metrics are exposed at:
`http://localhost:8080/metrics`

2. Promtail ships JSON logs from:
`./logs/app.log`
to Loki (inside Docker).

3. Grafana is available at the default port:
`http://localhost:3000`

Notes:
- If Prometheus can’t reach your locally running API, update `observability/prometheus.yml` target from `host.docker.internal:8080`.

## 🧪 Testing

- Unit tests + integration tests:
`make test`

- Coverage report:
`make test-coverage` (outputs `coverage.html`)

Before running tests, ensure SQLC code is generated (the Makefile calls `make sqlc`).
