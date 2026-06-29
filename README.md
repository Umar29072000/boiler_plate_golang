# Fiber Boilerplate - Golang Backend

Backend boilerplate menggunakan **Golang**, **Fiber**, **PostgreSQL**, **GORM**, dan **Redis**. Boilerplate ini dirancang dengan clean architecture, authentication lengkap, email system, security hardening, caching, performance optimization, monitoring, logging, Docker support, dan standardized API response.

Project ini menggunakan CLI berbasis Cobra. REST server dijalankan melalui command `rest`.

---

## 🚀 Feature Overview

### 🔐 Authentication & Authorization
- JWT authentication dengan **single token**.
- Token expiry default: `168h` / 7 hari.
- Register dan login mengembalikan response `{ user, token }`.
- Password hashing menggunakan Bcrypt.
- Last login tracking.
- Auth middleware dengan `Authorization: Bearer <token>`.
- Admin middleware berbasis role `admin`.

### ✅ Request Validation
- Request validation menggunakan `ozzo-validation`.
- Request DTO berada di `internal/rest/request`.
- Validation dipanggil dari handler dengan pattern `req.Validate()`.
- Response invalid validation menggunakan message `INVALID_VALIDATION`.

### 📦 DTO & Domain Organization
- Domain model/response berada di `domains/`.
- DTO data flow service/repository berada di `domains/dto/`.

### 🧱 Architecture Pattern
- Fiber tetap digunakan sebagai HTTP framework.
- Pattern handler/service/repository mengikuti style `star-fruit`.

### 🪵 Logging
- Structured logging menggunakan Logrus.
- Setiap flow penting memakai tag.
- Middleware request logging tersedia.

### 🛡️ Security Middleware
- Security headers middleware.
- CORS configurable via env.
- Rate limiter dengan Redis dan fallback memory.
- Compression middleware.
- Recover middleware.
- Request ID middleware.

### 🗄️ Database & ORM
- PostgreSQL.
- GORM ORM.
- Auto migration optional via `DB_MIGRATE=yes`.
- User model dengan soft delete dan timestamps.

### ⚡ Redis, Cache, Email
- Redis client tersedia.
- Cache package tersedia di `pkg/cache`.
- SMTP email package dan HTML templates tersedia di `pkg/email`.

### 🐳 Docker
- Dockerfile dan Docker Compose tersedia.
- Docker Compose menyediakan PostgreSQL, Redis, dan API service.

---

## 📁 Project Structure

```text
boiler_plate_be_golang/
├── app/
│   ├── main.go
│   ├── cmd/
│   │   ├── root.go
│   │   └── rest.go
│   └── config/
│       ├── app.go
│       ├── email.go
│       ├── jwt.go
│       ├── postgres.go
│       ├── redis.go
│       └── root.go
├── domains/
│   ├── auth.go
│   ├── user.go
│   └── dto/
│       └── user.go
├── internal/
│   ├── database/
│   │   └── migrations/
│   │       └── migrate.go
│   ├── repository/
│   │   └── user.go
│   ├── rest/
│   │   ├── auth.go
│   │   ├── health.go
│   │   ├── user.go
│   │   ├── request/
│   │   │   ├── auth.go
│   │   │   ├── base.go
│   │   │   ├── pagination.go
│   │   │   └── user.go
│   │   └── response/
│   │       └── response.go
│   └── service/
│       ├── auth.go
│       └── user.go
├── middleware/
│   ├── auth.go
│   ├── compress.go
│   ├── ratelimit.go
│   ├── requestlog.go
│   └── security.go
├── pkg/
│   ├── cache/
│   ├── constant/
│   ├── email/
│   ├── logger/
│   ├── model/
│   ├── redis/
│   ├── utils/
│   └── validator/
├── .env.example
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
- Redis 7+ optional, direkomendasikan untuk rate limiter/cache
- Docker & Docker Compose optional

---

## ⚙️ Environment Variables

Copy `.env.example` menjadi `.env`:

```bash
cp .env.example .env
```

Konfigurasi utama:

```env
APP_NAME=Fiber Boilerplate
APP_ENV=development
APP_PORT=3000
APP_URL=http://localhost:3000

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=fiber_boilerplate
DB_SSL_MODE=disable
DB_TIMEZONE=Asia/Jakarta

JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=168h

CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

EMAIL_FROM=noreply@fiberboilerplate.com
EMAIL_HOST=smtp.ethereal.email
EMAIL_PORT=587
EMAIL_USERNAME=your-ethereal-username
EMAIL_PASSWORD=your-ethereal-password

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

RATE_LIMIT_MAX=100
RATE_LIMIT_DURATION=15m
```

> Set `DB_MIGRATE=yes` jika ingin menjalankan auto migration saat aplikasi start.

---

## 🚀 Running Locally

### 1. Install Dependencies

```bash
go mod download
```

### 2. Setup Database

```sql
CREATE DATABASE fiber_boilerplate;
```

### 3. Run REST Server

```bash
go run ./app rest --env .env
```

Server berjalan di:

```text
http://localhost:3000
```

Base API path:

```text
/api/v1
```

### 4. Build Binary

```bash
go build -o app.exe ./app
./app.exe rest --env .env
```

Windows PowerShell:

```powershell
.\app.exe rest --env .env
```

---

## 🐳 Running with Docker Compose

```bash
docker-compose up -d
```

Services:

| Service | URL/Port |
|---|---|
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

Semua endpoint berada di prefix:

```text
/api/v1
```

### Health

#### Basic Health

```http
GET /api/v1/health
```

Example response:

```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "status": "healthy"
  }
}
```

#### Readiness Probe

```http
GET /api/v1/health/ready
```

#### Liveness Probe

```http
GET /api/v1/health/live
```

---

### Authentication

#### Register

```http
POST /api/v1/auth/register
Content-Type: application/json
```

Request:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "Password123"
}
```

Success response:

```json
{
  "code": 201,
  "message": "USER_REGISTERED",
  "data": {
    "user": {
      "id": "uuid",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "isEmailVerified": false
    },
    "token": "jwt-token"
  }
}
```

#### Login

```http
POST /api/v1/auth/login
Content-Type: application/json
```

Request:

```json
{
  "email": "john@example.com",
  "password": "Password123"
}
```

Success response:

```json
{
  "code": 200,
  "message": "LOGIN_SUCCESS",
  "data": {
    "user": {
      "id": "uuid",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user"
    },
    "token": "jwt-token"
  }
}
```

#### Verify Email

```http
GET /api/v1/auth/verify/:token
```

#### Forgot Password

```http
POST /api/v1/auth/forgot-password
Content-Type: application/json
```

Request:

```json
{
  "email": "john@example.com"
}
```

#### Reset Password

```http
POST /api/v1/auth/reset-password
Content-Type: application/json
```

Request:

```json
{
  "token": "reset-token",
  "newPassword": "NewPassword123"
}
```

> Tidak ada endpoint `/auth/refresh` karena authentication menggunakan single token.

---

### User Management

Gunakan header:

```http
Authorization: Bearer <token>
```

#### Get Profile

```http
GET /api/v1/users/profile
```

#### Update Profile

```http
PUT /api/v1/users/profile
Content-Type: application/json
```

Request:

```json
{
  "name": "John Updated"
}
```

#### Get All Users

```http
GET /api/v1/users?page=1&limit=10&search=john&field=createdAt&sort=desc
```

#### Delete User

```http
DELETE /api/v1/users/:id
```

---

## 🧪 Testing with cURL

### Health

```bash
curl http://localhost:3000/api/v1/health
curl http://localhost:3000/api/v1/health/ready
curl http://localhost:3000/api/v1/health/live
```

### Register

```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"Password123"}'
```

### Login

```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Password123"}'
```

### Get Profile

```bash
curl -X GET http://localhost:3000/api/v1/users/profile \
  -H "Authorization: Bearer <token>"
```

### Update Profile

```bash
curl -X PUT http://localhost:3000/api/v1/users/profile \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Name"}'
```

### Get Users

```bash
curl -X GET "http://localhost:3000/api/v1/users?page=1&limit=10" \
  -H "Authorization: Bearer <token>"
```

---

## 📤 Standard Response

### Base Response

```json
{
  "code": 200,
  "message": "SUCCESS",
  "data": {}
}
```

### Pagination Data

```json
{
  "code": 200,
  "message": "SUCCESS",
  "data": {
    "items": [],
    "totalItems": 0,
    "totalPages": 0,
    "page": 1,
    "limit": 10
  }
}
```

### Validation Error

```json
{
  "code": 400,
  "message": "INVALID_VALIDATION",
  "data": {
    "field": "error message"
  }
}
```

---

## 🪵 Logging Pattern

Pattern logging mengikuti style tag bertingkat:

```go
var tag string = "internal.rest.auth.Login."

logrus.WithFields(logrus.Fields{
    "tag":   tag + "01",
    "error": err.Error(),
}).Error("bad request")
```

Contoh tag:

| Layer | Example Tag |
|---|---|
| REST | `internal.rest.auth.Login.01` |
| Service | `internal.service.auth.Login.03` |
| Repository | `internal.repository.user.Show.01` |

---

## 🧭 Development Guide

### Add New Feature Module

1. Tambahkan domain/entity di `domains/`.
2. Tambahkan DTO flow di `domains/dto/` jika dibutuhkan.
3. Tambahkan request struct dan `Validate()` di `internal/rest/request/`.
4. Tambahkan response struct khusus di `internal/rest/response/` jika dibutuhkan.
5. Buat repository di `internal/repository/`.
6. Buat service di `internal/service/`.
7. Buat handler di `internal/rest/`.
8. Register handler di `app/cmd/rest.go`.

### Handler Pattern

```go
type exampleHandler struct {
    ExampleService service.IExampleService
}

func (h *exampleHandler) Create(c *fiber.Ctx) error {
    var (
        tag string = "internal.rest.example.Create."
        req request.CreateExampleRequest
    )

    if err := c.BodyParser(&req); err != nil {
        logrus.WithFields(logrus.Fields{
            "tag":   tag + "01",
            "error": err.Error(),
        }).Error("bad request")
        return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
            Code:    http.StatusBadRequest,
            Message: "BAD_REQUEST",
        })
    }

    if err := req.Validate(); err != nil {
        return c.Status(http.StatusBadRequest).JSON(response.BaseResponse{
            Code:    http.StatusBadRequest,
            Message: "INVALID_VALIDATION",
            Data:    err,
        })
    }

    return c.Status(http.StatusOK).JSON(response.BaseResponse{
        Code:    http.StatusOK,
        Message: "SUCCESS",
    })
}
```

### Protected Route

```go
users := api.Group("/users")
users.Use(middleware.Auth())
users.Get("/profile", handler.GetProfile)
```

### Admin Route

```go
users.Get("/", middleware.AdminOnly(), handler.GetAllUsers)
```

---

## 📦 Main Dependencies

- [Fiber](https://github.com/gofiber/fiber)
- [Cobra](https://github.com/spf13/cobra)
- [GORM](https://gorm.io/)
- [PostgreSQL Driver](https://gorm.io/docs/connecting_to_the_database.html)
- [go-redis](https://github.com/redis/go-redis)
- [JWT](https://github.com/golang-jwt/jwt)
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [Logrus](https://github.com/sirupsen/logrus)
- [Zerolog](https://github.com/rs/zerolog)
- [Ozzo Validation](https://github.com/go-ozzo/ozzo-validation)
- [godotenv](https://github.com/joho/godotenv)

---

## ✅ Verification

Build aplikasi:

```bash
go build -o app.exe ./app
```

Run server:

```bash
./app.exe rest --env .env
```

---

## 📄 License

MIT License
