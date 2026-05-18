# E-Commerce REST API

A production-structured REST API built with **Go**, **Gin**, **GORM**, and **PostgreSQL**.
Demonstrates Clean Architecture with strict layer separation — Handler → UseCase → Repository → Entity.

Built as a 30-day learning project to master backend development in Go.

## Tech Stack

| Technology   | Version | Purpose                        |
|-------------|---------|--------------------------------|
| Go           | 1.22    | Language                       |
| Gin          | v1.9    | HTTP framework                 |
| GORM         | v1.25   | ORM + PostgreSQL driver        |
| PostgreSQL   | 15      | Database                       |
| JWT          | v5      | Authentication tokens          |
| bcrypt       | stdlib  | Password hashing               |
| Docker       | latest  | Containerisation               |

## Architecture

This project follows Clean Architecture — dependencies point inward only.

```text
Handler (Gin)  →  UseCase (Business Logic)  →  Repository (GORM)  →  Entity (Pure Go)
```

| Layer      | Location                      | Responsibility                          |
|-----------|-------------------------------|----------------------------------------|
| Entity     | internal/domain/entity/       | Pure structs, business rules, interfaces|
| UseCase    | internal/usecase/             | Business logic, no framework imports    |
| Repository | internal/repository/postgres/ | GORM queries, model↔entity mapping      |
| Handler    | internal/delivery/http/       | HTTP binding, response formatting       |

**The test:** `grep -r "gin-gonic\|gorm.io" internal/domain/ internal/usecase/`
returns nothing — the inner layers are completely framework-free.

## Quick Start

### With Docker (recommended)

```bash
git clone https://github.com/YOUR_USERNAME/ecommerce-api.git
cd ecommerce-api
cp .env.example .env
docker-compose up --build
```

API is running at `http://localhost:8080`

### Local development

```bash
# Prerequisites: Go 1.22+, PostgreSQL 15
git clone https://github.com/YOUR_USERNAME/ecommerce-api.git
cd ecommerce-api
cp .env.example .env        # edit DATABASE_URL
go mod tidy
go run cmd/api/main.go
```

## Environment Variables

Copy `.env.example` to `.env` and update values.

| Variable       | Default                              | Description              |
|---------------|--------------------------------------|--------------------------|
| `PORT`        | `8080`                               | Server port              |
| `DATABASE_URL`| `postgres://...@localhost:5432/...`  | PostgreSQL connection    |
| `JWT_SECRET`  | —                                    | JWT signing secret (required) |

## API Endpoints

All responses follow: `{"success": bool, "data": ..., "error": "..."}`

### Auth (public)

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Rahim","email":"rahim@example.com","password":"securepass123"}'

# Login → returns JWT token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"rahim@example.com","password":"securepass123"}'
```

### Products (read: public / write: requires Bearer token)

```bash
# List products (with search, filter, sort, pagination)
curl "http://localhost:8080/api/v1/products?search=laptop&sort=price_asc&page=1&limit=10"

# Get one product
curl http://localhost:8080/api/v1/products/1

# Create product (requires auth)
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Laptop Pro","description":"Fast laptop","price":999.99,"stock":10}'

# Update product (requires auth)
curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Laptop Pro Max","price":1099.99,"stock":8}'

# Delete product (requires auth)
curl -X DELETE http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Orders (requires Bearer token)

```bash
# Place an order (stock checked + price snapshot)
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"items":[{"product_id":1,"quantity":2}]}'

# Get order by ID (owner only)
curl http://localhost:8080/api/v1/orders/1 \
  -H "Authorization: Bearer $TOKEN"

# List my orders (paginated)
curl "http://localhost:8080/api/v1/orders/my?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"

# Cancel order (pending orders only)
curl -X PATCH http://localhost:8080/api/v1/orders/1/cancel \
  -H "Authorization: Bearer $TOKEN"
```

### Health

```bash
curl http://localhost:8080/health
# {"success":true,"status":"ok","uptime":"5m3s"}
```

## Project Structure

```text
ecommerce-api/
├── cmd/api/main.go                    # Entry point — DI root
├── config/                            # Config loading + DB connection
├── internal/
│   ├── domain/
│   │   ├── entity/                    # Pure domain structs
|   |   ├── contracts/                 # Pure usecase interfaces
│   │   └── repository/                # Repository interface contracts
│   ├── usecase/                       # Business logic (no frameworks)
│   ├── repository/
│   │   ├── models/                    # GORM models (DB layer only)
│   │   └── postgres/                  # GORM implementations
│   └── delivery/http/
│       ├── handler/                   # Gin handlers
│       ├── middleware/                # Auth + Logger + Recovery
│       ├── response.go                # Unified response wrapper
│       └── router.go                  # Route registration
└── pkg/
    ├── apperror/                      # Sentinel errors
    ├── jwt/                           # Token generation + validation
    ├── hash/                          # bcrypt helpers
    └── pagination/                    # Shared pagination helper
```

**Dependency direction (enforced by Go's import system):**

```text
Handler → UseCase → Repository → Entity
```

No layer imports from a layer outside it. Circular imports are impossible.

## Key Features

- ✅ Full CRUD for Products with search, filtering, and sorting
- ✅ User registration + login with bcrypt password hashing
- ✅ JWT authentication with protected/public route groups
- ✅ Order placement with stock checking and price snapshots
- ✅ Paginated responses with total_pages metadata
- ✅ Consistent JSON response shape across all endpoints
- ✅ Field-level validation error detail
- ✅ Request logging with colour-coded status codes
- ✅ Panic recovery — server never crashes
- ✅ Docker + docker-compose — one command to run everything
- ✅ Multi-stage Docker build — ~20MB final image

## What I Learned

Building this project from zero Go knowledge taught me that Clean Architecture
is not about following rules — it is about making change cheap. When I replaced
the in-memory repository with PostgreSQL, only one file changed. That moment
made every architectural decision from the previous three weeks make sense.

The toughest concept was interfaces — specifically that the domain layer owns
its contracts, not the implementations. Once that clicked, Dependency Injection
became obvious rather than magical.

I also learned that error handling is not an afterthought. Go's explicit error
values forced me to think about every failure path, which produced more reliable
code than I was able to write in languages with exceptions.

## License

MIT
