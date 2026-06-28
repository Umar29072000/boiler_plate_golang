# Fiber Boilerplate - Golang Backend

Production-ready backend boilerplate menggunakan **Golang**, **Fiber**, **PostgreSQL**, **GORM**, dan **Redis**. Boilerplate ini dirancang dengan clean architecture, authentication lengkap, email system, security hardening, caching, performance optimization, monitoring, logging, Docker support, dan standardized API response.

## 🚀 Feature Overview

### 🔐 Authentication & Authorization
- ✅ JWT Authentication menggunakan single access token
- ✅ Token expiry default: `7 days` (`168h`)
- ✅ JWT verification middleware
- ✅ Password hashing menggunakan Bcrypt
- ✅ Password comparison helper
- ✅ Last login tracking
- ✅ Role-Based Access Control (RBAC)
- ✅ Role tersedia: `user`, `admin`
- ✅ Auth/protect middleware
- ✅ Admin authorization middleware
- ✅ Route-level protection

### 📧 Email System
- ✅ SMTP email service
- ✅ Development email testing via Ethereal Email
- ✅ Production SMTP support: Gmail, SendGrid, Mailgun, dan SMTP lainnya
- ✅ Welcome email setelah registration
- ✅ Email verification token dengan expiry 24 jam
- ✅ Password reset token dengan expiry 15 menit
- ✅ Password changed confirmation email
- ✅ Resend verification email
- ✅ Async email sending
- ✅ Professional HTML templates:
  - `welcome.html`
  - `verifyEmail.html`
  - `resetPassword.html`
  - `passwordChanged.html`

### 🛡️ Security Enhancement
- ✅ Security headers middleware, setara Helmet.js
- ✅ XSS protection headers
- ✅ Clickjacking protection
- ✅ MIME sniffing protection
- ✅ DNS prefetch control
- ✅ Referrer policy
- ✅ Permissions policy
- ✅ CORS configurable by environment
- ✅ Redis-backed distributed rate limiting
- ✅ Memory fallback jika Redis tidak tersedia
- ✅ Default rate limit: 100 requests per 15 minutes
- ✅ Rate limit response headers
- ✅ XSS input protection
- ✅ Input sanitization
- ✅ SQL injection pattern defense-in-depth
- ✅ `.env` excluded by `.gitignore`

### ⚡ Caching & Performance
- ✅ Redis cache integration
- ✅ High-level cache service
- ✅ Response caching middleware for GET requests
- ✅ Configurable cache TTL
- ✅ Cache hit/miss response header: `X-Cache`
- ✅ Cache hit/miss statistics support
- ✅ Automatic cache invalidation middleware
- ✅ Pattern-based cache deletion
- ✅ Cache warm-up helper
- ✅ Cache operations:
  - `Get`
  - `Set`
  - `Del`
  - `DelPattern`
  - `Exists`
  - `Flush`
  - `GetStats`
  - `InvalidateByPrefix`
  - `WarmCache`
- ✅ Gzip response compression
- ✅ Compression level options: default, best compression, best speed

### 🏥 Monitoring & Logging
- ✅ Enhanced health checks
- ✅ Basic health endpoint
- ✅ Detailed health endpoint
- ✅ Kubernetes-style readiness probe
- ✅ Kubernetes-style liveness probe
- ✅ DB connection status monitoring
- ✅ Redis connection status monitoring
- ✅ Server uptime tracking
- ✅ Runtime metrics endpoint
- ✅ Go runtime memory metrics
- ✅ Goroutine count
- ✅ GC metrics
- ✅ Structured logging with Zerolog
- ✅ Pretty console logs in development
- ✅ JSON logs in production
- ✅ Request ID middleware
- ✅ `X-Request-ID` response header
- ✅ Structured request/response logging
- ✅ Logs method, path, IP, user-agent, status, duration, response size, and errors

### 🗄️ Database & ORM
- ✅ PostgreSQL database
- ✅ GORM ORM
- ✅ Auto migration
- ✅ Graceful DB close
- ✅ Database connection status checks
- ✅ Soft delete support via GORM
- ✅ Automatic timestamps: `created_at`, `updated_at`
- ✅ User model with:
  - Email uniqueness
  - Role
  - Email verification status
  - Email verification token
  - Password reset token
  - Last login timestamp

### 🐳 Docker & DevOps
- ✅ Multi-stage Dockerfile
- ✅ Docker Compose support
- ✅ 3-service orchestration:
  - API
  - PostgreSQL
  - Redis
- ✅ PostgreSQL data persistence
- ✅ Redis data persistence
- ✅ Health checks for PostgreSQL and Redis
- ✅ Automatic service dependency health checks
- ✅ Restart policies
- ✅ Environment-based configuration

### 🧱 Architecture
- ✅ Clean Architecture inspired structure
- ✅ Separation of concerns:
  - Controllers
  - Services
  - Repositories
  - Models
  - Middleware
  - Utilities
- ✅ Standardized API responses
- ✅ Centralized error handling
- ✅ Reusable utilities

---

## 📁 Project Structure

```text
boiler_plate_be_golang/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── controllers/
│   │   ├── auth.go
│   │   ├── health.go
│   │   ├── metrics.go
│   │   └── user.go
│   ├── database/
│   │   ├── database.go
│   │   └── migrations/
│   │       └── migrate.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cacheresp.go
│   │   ├── compress.go
│   │   ├── cors.go
│   │   ├── error.go
│   │   ├── logger.go
│   │   ├── ratelimit.go
│   │   ├── requestlog.go
│   │   ├── sanitize.go
│   │   └── security.go
│   ├── models/
│   │   ├── base.go
│   │   └── user.go
│   ├── repositories/
│   │   └── user.go
│   ├── routes/
│   │   └── routes.go
│   └── services/
│       ├── auth.go
│       └── user.go
├── pkg/
│   ├── cache/
│   │   └── cache.go
│   ├── email/
│   │   ├── email.go
│   │   └── templates/
│   │       ├── passwordChanged.html
│   │       ├── resetPassword.html
│   │       ├── verifyEmail.html
│   │       └── welcome.html
│   ├── logger/
│   │   └── logger.go
│   ├── redis/
│   │   └── redis.go
│   ├── utils/
│   │   ├── jwt.go
│   │   ├── password.go
│   │   ├── response.go
│   │   └── token.go
│   └── validator/
│       └── validator.go
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

## 🛠️ Prerequisites

- Go 1.22+
- PostgreSQL 14+
- Redis 7+ optional, tapi direkomendasikan
- Docker & Docker Compose optional

---

## ⚙️ Environment Variables

Copy file `.env.example` menjadi `.env`:

```bash
cp .env.example .env
```

Konfigurasi utama:

```env
# Application
APP_NAME=Fiber Boilerplate
APP_ENV=development
APP_PORT=3000
APP_URL=http://localhost:3000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=fiber_boilerplate
DB_SSL_MODE=disable
DB_TIMEZONE=Asia/Jakarta

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=168h

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# Email
EMAIL_FROM=noreply@fiberboilerplate.com
EMAIL_HOST=smtp.ethereal.email
EMAIL_PORT=587
EMAIL_USERNAME=your-ethereal-username
EMAIL_PASSWORD=your-ethereal-password

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Rate Limiting
RATE_LIMIT_MAX=100
RATE_LIMIT_DURATION=15m
```

---

## 🚀 Running Locally

### 1. Install Dependencies

```bash
go mod download
```

### 2. Setup Database

Buat database PostgreSQL:

```sql
CREATE DATABASE fiber_boilerplate;
```

### 3. Run App

```bash
go run cmd/api/main.go
```

Server berjalan di:

```text
http://localhost:3000
```

### 4. Build Binary

```bash
go build -o main.exe cmd/api/main.go
./main.exe
```

---

## 🐳 Running with Docker Compose

```bash
docker-compose up -d
```

Services:

| Service | URL/Port |
|---------|----------|
| API | `http://localhost:3000` |
| PostgreSQL | `localhost:5432` |
| Redis | `localhost:6379` |

View logs:

```bash
docker-compose logs -f
```

Stop services:

```bash
docker-compose down
```

Stop and remove volumes:

```bash
docker-compose down -v
```

---

## 📡 API Endpoints

### Health & Monitoring

#### Basic Health

```http
GET /api/health
```

#### Detailed Health

```http
GET /api/health/detailed
```

Example response:

```json
{
  "status": "OK",
  "timestamp": "2026-06-28T17:00:00Z",
  "uptime": "10m15s",
  "environment": "development",
  "version": "Fiber Boilerplate",
  "services": {
    "database": {
      "status": "healthy",
      "type": "PostgreSQL"
    },
    "redis": {
      "status": "healthy",
      "type": "Redis"
    }
  }
}
```

#### Readiness Probe

```http
GET /api/health/ready
```

#### Liveness Probe

```http
GET /api/health/live
```

#### Metrics

```http
GET /api/metrics
```

---

### Authentication

#### Register

```http
POST /api/auth/register
Content-Type: application/json
```

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

#### Login

```http
POST /api/auth/login
Content-Type: application/json
```

```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

#### Verify Email

```http
GET /api/auth/verify-email/:token
```

#### Resend Verification Email

```http
POST /api/auth/resend-verification
Content-Type: application/json
```

```json
{
  "email": "john@example.com"
}
```

#### Forgot Password

```http
POST /api/auth/forgot-password
Content-Type: application/json
```

```json
{
  "email": "john@example.com"
}
```

#### Reset Password

```http
POST /api/auth/reset-password/:token
Content-Type: application/json
```

```json
{
  "password": "newPassword123"
}
```

---

### User Management Protected

Gunakan header:

```http
Authorization: Bearer <token>
```

#### Get Profile

```http
GET /api/users/profile
```

#### Update Profile

```http
PUT /api/users/profile
Content-Type: application/json
```

```json
{
  "name": "John Updated"
}
```

#### Get All Users Admin Only

```http
GET /api/users?page=1&limit=10
```

#### Delete User Admin Only

```http
DELETE /api/users/:id
```

---

## 🧪 Testing with cURL

### Health

```bash
curl http://localhost:3000/api/health
curl http://localhost:3000/api/health/detailed
curl http://localhost:3000/api/metrics
```

### Register

```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'
```

### Login

```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Get Profile

```bash
curl -X GET http://localhost:3000/api/users/profile \
  -H "Authorization: Bearer <token>"
```

### Test Security Headers

```bash
curl -I http://localhost:3000/api/health
```

### Test Compression

```bash
curl -H "Accept-Encoding: gzip" -I http://localhost:3000/api/health
```

### Test Rate Limit

```bash
for i in {1..101}; do curl http://localhost:3000/api/health; done
```

---

## 📧 Email Setup

### Development with Ethereal

1. Buka https://ethereal.email
2. Buat test account
3. Masukkan credential ke `.env`

```env
EMAIL_HOST=smtp.ethereal.email
EMAIL_PORT=587
EMAIL_USERNAME=your-ethereal-username@ethereal.email
EMAIL_PASSWORD=your-ethereal-password
```

### Production with Gmail

```env
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-app-password
```

### Production with SendGrid

```env
EMAIL_HOST=smtp.sendgrid.net
EMAIL_PORT=587
EMAIL_USERNAME=apikey
EMAIL_PASSWORD=your-sendgrid-api-key
```

---

## ⚡ Cache Usage Example

```go
api.Get("/products",
    middleware.Cache(middleware.CacheConfig{
        TTL: 10 * time.Minute,
    }),
    getProductsHandler,
)

api.Post("/products",
    middleware.Auth(),
    middleware.InvalidateCache("GET:/api/products"),
    createProductHandler,
)
```

Manual cache service:

```go
cacheService := cache.NewCacheService(5 * time.Minute)
cacheService.Set("user:1", user, 10*time.Minute)

var cachedUser User
cacheService.Get("user:1", &cachedUser)
```

---

## 🪵 Logging

Development mode menghasilkan pretty console logs.

Production mode menghasilkan JSON logs yang cocok untuk log aggregator seperti:

- ELK Stack
- Grafana Loki
- Datadog
- CloudWatch
- Google Cloud Logging

Contoh structured log:

```json
{
  "level": "info",
  "request_id": "uuid",
  "method": "GET",
  "path": "/api/health",
  "status": 200,
  "duration": 1200000,
  "message": "Request completed"
}
```

---

## 🔐 Security Headers

Middleware security menambahkan header berikut:

- `X-XSS-Protection`
- `X-Content-Type-Options`
- `X-Frame-Options`
- `X-DNS-Prefetch-Control`
- `Referrer-Policy`
- `Permissions-Policy`

---

## 📦 Main Dependencies

- [Fiber](https://github.com/gofiber/fiber)
- [GORM](https://gorm.io/)
- [PostgreSQL Driver](https://gorm.io/docs/connecting_to_the_database.html)
- [go-redis](https://github.com/redis/go-redis)
- [JWT](https://github.com/golang-jwt/jwt)
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [godotenv](https://github.com/joho/godotenv)
- [zerolog](https://github.com/rs/zerolog)

---

## 🧭 Development Guide

### Add New Model

1. Buat model di `internal/models/`
2. Tambahkan model ke `internal/database/migrations/migrate.go`
3. Buat repository di `internal/repositories/`
4. Buat service di `internal/services/`
5. Buat controller di `internal/controllers/`
6. Register route di `internal/routes/routes.go`

### Add Protected Route

```go
users := api.Group("/users")
users.Use(middleware.Auth())
users.Get("/profile", userController.GetProfile)
```

### Add Admin Route

```go
users.Get("/", middleware.AdminOnly(), userController.GetAllUsers)
```

---

## ✅ Roadmap Status

- ✅ Phase 1: Email System
- ✅ Phase 2: Security Enhancement
- ✅ Phase 3: Caching & Performance
- ✅ Phase 4: Monitoring & Logging

---

## 📄 License

MIT License
