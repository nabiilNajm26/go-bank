# GoBank

Banking API built with Go and PostgreSQL.

## Tech Stack

- Go 1.21
- Fiber v2
- PostgreSQL
- Redis
- Docker
- AWS S3 & ECS

## Setup

```bash
# Install dependencies
go mod download

# Run migrations
make migrate-up

# Start server
make run
```

## API Endpoints

- `POST /auth/register`
- `POST /auth/login`
- `GET /accounts`
- `POST /transfers`
- `GET /transactions`

## Project Structure

```
cmd/api/        - Main application
internal/       - Business logic
db/migrations/  - SQL migrations
deployments/    - Docker & K8s configs
```

## Testing

```bash
make test
```

## Docker

```bash
docker-compose up
```

## Author

Nabiil Najm