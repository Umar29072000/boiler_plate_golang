# Fiber Boilerplate - Golang Backend

Boilerplate backend API menggunakan Golang dengan Fiber framework, siap pakai untuk membangun RESTful API dengan fitur authentication, authorization, dan clean architecture.

## 🚀 Fitur

- ✅ **Fiber Framework** - Web framework Go yang cepat dan ekspresif
- ✅ **Clean Architecture** - Struktur kode yang terorganisir (Controllers, Services, Repositories)
- ✅ **JWT Authentication** - Authentication menggunakan JSON Web Token (7 days expiry)
- ✅ **Email System** - Complete email functionality with SMTP support
- ✅ **Role-based Access Control** - Authorization dengan role (user, admin)
- ✅ **PostgreSQL Database** - Database relational dengan GORM ORM
- ✅ **Auto Migration** - Database migration otomatis
- ✅ **Last Login Tracking** - Track user login activity
- ✅ **Middleware** - CORS, Logger, Error Handler, Auth
- ✅ **Input Validation** - Validasi input request
- ✅ **Password Hashing** - Bcrypt password hashing
- ✅ **Environment Config** - Konfigurasi melalui environment variables
- ✅ **Docker Support** - Docker dan Docker Compose ready
- ✅ **Standardized Response** - Response format yang konsisten

## 📁 Struktur Folder

```
boiler_plate_be_golang/
├── cmd/
│   └── api/
│       └── main.go              # Entry point aplikasi
├── internal/
│   ├── config/
│   │   └── config.go            # Konfigurasi aplikasi
│   ├── controllers/
│   │   ├── auth.go              # Auth controller
│   │   └── user.go              # User controller
│   ├── database/
│   │   ├── database.go          # Database connection
│   │   └── migrations/
│   │       └── migrate.go       # Auto migration
│   ├── middleware/
│   │   ├── auth.go              # JWT middleware
│   │   ├── cors.go              # CORS middleware
│   │   ├── error.go             # Error handler
│   │   └── logger.go            # Request logger
│   ├── models/
│   │   ├── base.go              # Base model
│   │   └── user.go              # User model
│   ├── repositories/
│   │   └── user.go              # User repository
│   ├── routes/
│   │   └── routes.go            # Route definitions
│   └── services/
│       ├── auth.go              # Auth business logic
│       └── user.go              # User business logic
├── pkg/
│   ├── utils/
│   │   ├── jwt.go               # JWT utilities
│   │   ├── password.go          # Password utilities
│   │   └── response.go          # Response formatter
│   └── validator/
│       └── validator.go         # Input validator
├── .env.example                 # Environment variables template
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

## 🛠️ Instalasi

### Prerequisites

- Go 1.22 atau lebih tinggi
- PostgreSQL 14+
- Docker & Docker Compose (optional)

### Setup Manual

1. **Clone repository**
```bash
git clone <repository-url>
cd boiler_plate_be_golang
```

2. **Install dependencies**
```bash
go mod download
```

3. **Setup environment variables**
```bash
cp .env.example .env
```

Edit file `.env` sesuai konfigurasi Anda:
```env
APP_NAME=Fiber Boilerplate
APP_ENV=development
APP_PORT=3000

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=fiber_boilerplate

JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRATION=24h
```

4. **Setup database**

Buat database PostgreSQL:
```sql
CREATE DATABASE fiber_boilerplate;
```

5. **Run aplikasi**
```bash
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:3000`

### Setup dengan Docker

1. **Build dan jalankan dengan Docker Compose**
```bash
docker-compose up -d
```

2. **Stop containers**
```bash
docker-compose down
```

3. **View logs**
```bash
docker-compose logs -f api
```

## 📡 API Endpoints

### Health Check
```
GET /api/health
```

### Authentication

#### Register
```
POST /api/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

#### Login
```
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "isEmailVerified": true,
      "lastLogin": "2024-01-15T10:30:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

#### Verify Email
```
GET /api/auth/verify-email/:token
```

Response:
```json
{
  "success": true,
  "message": "Email verified successfully",
  "data": null
}
```

#### Resend Verification Email
```
POST /api/auth/resend-verification
Content-Type: application/json

{
  "email": "john@example.com"
}
```

#### Forgot Password
```
POST /api/auth/forgot-password
Content-Type: application/json

{
  "email": "john@example.com"
}
```

Response:
```json
{
  "success": true,
  "message": "If your email exists in our system, you will receive a password reset link",
  "data": null
}
```

#### Reset Password
```
POST /api/auth/reset-password/:token
Content-Type: application/json

{
  "password": "newPassword123"
}
```

Response:
```json
{
  "success": true,
  "message": "Password reset successfully",
  "data": null
}
```

### User Management (Protected)

Semua endpoint user memerlukan JWT token di header:
```
Authorization: Bearer <token>
```

#### Get Profile
```
GET /api/users/profile
Authorization: Bearer <token>
```

#### Update Profile
```
PUT /api/users/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "John Updated"
}
```

#### Get All Users (Admin Only)
```
GET /api/users?page=1&limit=10
Authorization: Bearer <token>
```

#### Delete User (Admin Only)
```
DELETE /api/users/:id
Authorization: Bearer <token>
```

## 🔐 Authentication & Authorization

### JWT Token

Setelah login/register, Anda akan menerima JWT token. Gunakan token ini untuk mengakses protected endpoints dengan menambahkan ke header:

```
Authorization: Bearer <your-token>
```

### Roles

- **user** - Role default untuk user baru
- **admin** - Role dengan akses penuh

Untuk membuat admin, ubah role di database:
```sql
UPDATE users SET role = 'admin' WHERE email = 'admin@example.com';
```

## 🔧 Konfigurasi

Konfigurasi aplikasi melalui environment variables di file `.env`:

| Variable | Description | Default |
|----------|-------------|---------|
| APP_NAME | Nama aplikasi | Fiber Boilerplate |
| APP_ENV | Environment (development/production) | development |
| APP_PORT | Port server | 3000 |
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | postgres |
| DB_NAME | Database name | fiber_boilerplate |
| JWT_SECRET | JWT secret key | your-secret-key |
| JWT_EXPIRATION | Token expiration duration | 168h (7 days) |
| CORS_ALLOWED_ORIGINS | Allowed origins for CORS | * |
| EMAIL_FROM | Sender email address | noreply@fiberboilerplate.com |
| EMAIL_HOST | SMTP host | smtp.ethereal.email |
| EMAIL_PORT | SMTP port | 587 |
| EMAIL_USERNAME | SMTP username | (empty) |
| EMAIL_PASSWORD | SMTP password | (empty) |

### 📧 Email Configuration

#### Development (Ethereal Email)
For testing emails in development, use [Ethereal Email](https://ethereal.email):
1. Visit https://ethereal.email
2. Create a test account
3. Copy credentials to `.env`:
```env
EMAIL_HOST=smtp.ethereal.email
EMAIL_PORT=587
EMAIL_USERNAME=your-ethereal-username@ethereal.email
EMAIL_PASSWORD=your-ethereal-password
```

#### Production (Gmail)
```env
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-app-password
```

#### Production (SendGrid)
```env
EMAIL_HOST=smtp.sendgrid.net
EMAIL_PORT=587
EMAIL_USERNAME=apikey
EMAIL_PASSWORD=your-sendgrid-api-key
```

## 🧪 Testing dengan cURL

### Register User
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
  -H "Authorization: Bearer <your-token>"
```

## 📦 Dependencies

- [Fiber](https://github.com/gofiber/fiber) - Web framework
- [GORM](https://gorm.io/) - ORM library
- [JWT](https://github.com/golang-jwt/jwt) - JWT implementation
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - Password hashing
- [godotenv](https://github.com/joho/godotenv) - Environment variables loader

## 🚀 Production Deployment

### Build Binary
```bash
go build -o main cmd/api/main.go
```

### Run Binary
```bash
./main
```

### Environment
Pastikan set environment variables di production:
- Set `APP_ENV=production`
- Gunakan JWT secret yang kuat dan unik
- Konfigurasi database production
- Set CORS allowed origins sesuai frontend domain

## 📝 Development

### Add New Model
1. Buat file model di `internal/models/`
2. Tambahkan ke migration di `internal/database/migrations/migrate.go`

### Add New Endpoint
1. Buat repository di `internal/repositories/`
2. Buat service di `internal/services/`
3. Buat controller di `internal/controllers/`
4. Register route di `internal/routes/routes.go`

## 📄 License

MIT License

## 👥 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
