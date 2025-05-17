# Coupon System MVP (GoKit + PostgreSQL)

This is a Coupon System MVP built using **GoKit** and **PostgreSQL**, designed for a medicine ordering platform.

---

## ✨ Features

- Admin coupon creation with constraints
- Coupon validation (for cart, time, categories, usage type)
- Support for:
  - one-time / multi-use / time-based
  - medicine/category-based applicability
  - order value and time window conditions
- Applicable coupons API
- Concurrency-aware coupon usage
- In-memory caching
- Swagger/OpenAPI integration
- Dockerized app

---

## 📁 Project Structure

| Folder        | Purpose                              |
| ------------- | ------------------------------------ |
| `cmd/`        | Main application entry               |
| `coupon/`     | Business logic, endpoints, transport |
| `db/`         | DB schema scripts                    |
| `pkg/`        | Caching utilities                    |

---

## 🚀 Getting Started

### Prerequisites
- Go >= 1.20
- PostgreSQL >= 13
- Docker (optional)

### Setup Instructions

```bash
# Clone repo
git clone https://coupon-system.git
cd coupon-system

# Run DB
psql -U postgres -f db/schema.sql

# Build & run
go mod tidy
go run cmd/couponservice/main.go




coupon-system/
├── cmd/
│   └── couponservice/
│       └── main.go
├── coupon/
│   ├── service.go
│   ├── endpoint.go
│   ├── transport.go
│   ├── repository.go
│   ├── models.go
│   ├── validation.go
├── db/
│   └── schema.sql
├── pkg/
│   └── cache.go
├── go.mod
├── go.sum
├── Dockerfile
└── README.md

module github.com/yourusername/coupon-system

go 1.20

require (
    github.com/go-kit/kit/v0.12.0
    github.com/gorilla/mux v1.8.1
    github.com/lib/pq v1.10.9
    github.com/patrickmn/go-cache v2.1.0+incompatible
)

