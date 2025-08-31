# GoBank

Banking API built with Go and PostgreSQL. Created as a personal project to learn modern backend development.

## What's Inside

- Go with Fiber web framework
- PostgreSQL database
- JWT authentication
- Money transfers between accounts
- Rate limiting and security
- PDF/CSV statement generation
- Real-time notifications
- Docker support

## Getting Started

```bash
git clone https://github.com/nabiilNajm26/go-bank.git
cd go-bank

# Copy environment file
cp .env.example .env

# Start database (using Docker)
docker run -d --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=gobank \
  -p 5432:5432 postgres:15

# Run migrations
cat db/migrations/*.up.sql | psql -U postgres -d gobank

# Start the app
go run cmd/api/main.go
```

Server runs on `http://localhost:8080`

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
cmd/api/                 - Main application
internal/domain/         - Business models
internal/usecase/        - Business logic  
internal/delivery/http/  - HTTP handlers
internal/repository/     - Database access
db/migrations/           - SQL migrations
```

## Development

```bash
make run        # Start server
make test       # Run tests
make docker-up  # Start with Docker
```

## Features Completed

- User registration/login with JWT
- Account creation and management
- Money transfers with proper transactions
- Transaction history
- Rate limiting (security)
- PDF and CSV statement generation
- WebSocket real-time notifications
- Input validation
- Idempotency for transfers

## CI/CD

GitHub Actions runs tests and builds Docker images automatically. To push to Docker Hub, add these repository secrets:
- `DOCKER_USERNAME`
- `DOCKER_PASSWORD`

## Author

Built by Nabiil Najm as a learning project for modern Go development.