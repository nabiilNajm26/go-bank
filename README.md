# GoBank - Banking System

A production-ready banking API built with Go, featuring real AWS cloud integration and enterprise-grade architecture.

[![Go](https://img.shields.io/badge/Go-1.23-blue)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue)](https://postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-red)](https://redis.io/)
[![AWS S3](https://img.shields.io/badge/AWS-S3-orange)](https://aws.amazon.com/s3/)

## Production Features

This banking system includes real AWS S3 integration for file uploads, comprehensive caching with Redis, and production-ready security features.

```bash
# Test the API
curl http://localhost:8080/health

# View interactive documentation  
open http://localhost:8080/swagger/
```

## Architecture

Built with clean architecture principles and modern Go practices:

- **Clean Architecture** with domain-driven design
- **36 HTTP endpoints** with complete CRUD operations  
- **ACID transactions** for secure money transfers
- **Redis caching** with graceful database fallback
- **JWT authentication** with session management

## Security Features

- Rate limiting (100 requests/minute, 5 for auth endpoints)
- Idempotency middleware for financial transactions
- Input validation and sanitization  
- SQL injection prevention
- CORS configuration and security headers

## Cloud Integration

- **AWS S3** for profile image storage
- Automated cost controls and budget alerts
- File lifecycle policies for cost optimization
- Production-ready with proper error handling

## Getting Started

```bash
git clone https://github.com/nabiilNajm26/go-bank.git
cd go-bank

# Setup environment  
cp .env.example .env

# Start dependencies
docker compose up -d postgres redis

# Start the server
go run cmd/api/main.go
```

Access points:
- Server: http://localhost:8080  
- API Documentation: http://localhost:8080/swagger/
- Health Check: http://localhost:8080/health

## API Examples

Profile image upload with S3:
```bash
curl -X POST localhost:8080/api/v1/users/profile/image \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "image=@profile.png"
```

Money transfer between accounts:
```bash
curl -X POST localhost:8080/api/v1/transactions/transfer \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_account_id": "uuid",
    "to_account_id": "uuid", 
    "amount": "100.00"
  }'
```

## API Usage

Register a user:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "John Doe"
  }'
```

Create an account:
```bash
curl -X POST http://localhost:8080/api/v1/accounts \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "account_type": "savings",
    "currency": "USD"
  }'
```

## Project Structure

```
cmd/api/                    - Main application
internal/domain/            - Business models
internal/usecase/           - Business logic  
internal/delivery/http/     - HTTP handlers
internal/repository/        - Database access
internal/infrastructure/    - Redis, cache, sessions
db/migrations/              - SQL migrations
docs/                       - Swagger documentation
```

## Development

```bash
make run        # Start server
make test       # Run tests
make docker-up  # Start with Docker
```

## Complete Feature Set

### Core Banking Features
- User registration and authentication with JWT + Redis sessions  
- Account management with Redis caching
- Money transfers with ACID transaction support
- Transaction history with pagination and filtering
- PDF/CSV statement generation
- Real-time WebSocket notifications for account activities

### Security Implementation
- Rate limiting (100 requests/minute, 5/minute for auth)
- Idempotency middleware for financial transactions
- Comprehensive input validation and sanitization
- JWT-based authentication with session management
- SQL injection prevention and parameterized queries
- CORS configuration and security headers

### Infrastructure & DevOps
- AWS S3 integration for profile image storage
- Redis caching with database fallback
- Docker containerization for easy deployment
- GitHub Actions CI/CD pipeline  
- Interactive Swagger/OpenAPI documentation
- Structured logging and error handling

## Performance & Scalability

The system implements Redis caching with graceful database fallback:

```go
func (r *CachedUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
    // Try cache first (30min TTL)
    if user, err := r.cache.GetUser(ctx, id); err == nil {
        return user, nil
    }
    
    // Fallback to database
    user, err := r.userRepo.GetByID(ctx, id)
    if err == nil {
        r.cache.SetUser(ctx, user)
    }
    return user, err
}
```

## Project Structure

Built with clean architecture principles:

- **Domain Layer**: Business entities and rules
- **Use Case Layer**: Application business logic  
- **Infrastructure Layer**: Database, Redis, AWS S3
- **Delivery Layer**: HTTP handlers, middleware, WebSocket

```
cmd/api/                 - Main application entry point
internal/domain/         - Business models and entities
internal/usecase/        - Business logic layer
internal/delivery/http/  - HTTP handlers and middleware
internal/repository/     - Database access layer
internal/infrastructure/ - External services (Redis, S3)
db/migrations/           - Database schema migrations
docs/                    - API documentation
```

## Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| Backend | Go 1.23 + Fiber | High-performance HTTP server |
| Database | PostgreSQL 15 | ACID transactions and data persistence |
| Cache | Redis 7 | Session management and performance |
| Cloud Storage | AWS S3 | Profile image storage |
| Documentation | Swagger/OpenAPI | Interactive API documentation |
| DevOps | Docker + GitHub Actions | Containerization and CI/CD |

## Development

```bash
make run        # Start development server
make test       # Run test suite
make docker-up  # Start with Docker Compose
```

## Author

**Nabiil Najm** - Backend developer focused on scalable system design.

---

Star this repository if you find it helpful for learning Go web development or banking system architecture.