# SunReady API

A simplified solar project management API built with Go, featuring a clean architecture with repository and service patterns.

## Architecture

```
sunready/
├── cmd/
│   └── sunready/          # Application entry point
│       └── main.go
├── internal/
│   ├── database/          # Database connection
│   ├── models/            # Data models (GORM)
│   ├── repo/              # Repository layer (data access)
│   ├── service/           # Service layer (business logic)
│   ├── handler/           # HTTP handlers
│   └── middleware/        # HTTP middleware (auth, etc.)
├── db/
│   └── init.sql           # Database initialization
├── Dockerfile
├── docker-compose.yaml
└── go.mod
```

## Tech Stack

- **Language**: Go 1.20
- **Web Framework**: Chi Router
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Containerization**: Docker & Docker Compose

## Features

- ✅ User authentication (register/login with JWT)
- ✅ User management
- ✅ Company management
- ✅ Project management
- ✅ RESTful API design
- ✅ Clean architecture (repo/service/handler pattern)
- ✅ Docker support
- ✅ CORS enabled

## Prerequisites

- Go 1.20+ (for local development)
- Docker & Docker Compose (for containerized deployment)
- PostgreSQL 15+ (if running without Docker)

## Quick Start with Docker

1. **Clone and navigate to the project**:
   ```bash
   cd sunready
   ```

2. **Start the services**:
   ```bash
   docker-compose up -d
   ```

3. **Check the logs**:
   ```bash
   docker-compose logs -f api
   ```

4. **Test the API**:
   ```bash
   curl http://localhost:8080/health
   ```

The API will be available at `http://localhost:8080`

## Local Development Setup

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start PostgreSQL** (if not using Docker):
   ```bash
   # Using Docker for just the database
   docker run -d \
     --name sunready-postgres \
     -e POSTGRES_USER=sunready \
     -e POSTGRES_PASSWORD=sunready \
     -e POSTGRES_DB=sunready \
     -p 5432:5432 \
     postgres:15-alpine
   ```

4. **Initialize the database**:
   ```bash
   psql -h localhost -U sunready -d sunready -f db/init.sql
   ```

5. **Run the application**:
   ```bash
   go run cmd/sunready/main.go
   ```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `JWT_SECRET` | Secret key for JWT signing | Required |
| `PORT` | Server port | `8080` |
| `ENV` | Environment (development/production) | `development` |


## Building for Production

### Build binary
```bash
go build -o sunready ./cmd/sunready
```

### Build Docker image
```bash
docker build -t sunready:latest .
```

### Run with Docker
```bash
docker run -d \
  -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@host:5432/db" \
  -e JWT_SECRET="your-secret" \
  sunready:latest
```

## Development

### Run tests
```bash
go test ./...
```

### Format code
```bash
go fmt ./...
```

### Lint code
```bash
golangci-lint run
```

## Docker Commands

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f

# Rebuild and restart
docker-compose up -d --build

# Remove volumes (clean database)
docker-compose down -v
```

## Project Structure Explained

- **cmd/sunready**: Application entry point, wires up all dependencies
- **internal/models**: GORM models representing database tables
- **internal/repo**: Repository pattern - handles all database operations
- **internal/service**: Business logic layer - orchestrates repositories
- **internal/handler**: HTTP handlers - handle requests/responses
- **internal/middleware**: HTTP middleware (authentication, logging, etc.)
- **internal/database**: Database connection setup

## Security Notes

- Change `JWT_SECRET` in production
- Use strong passwords
- Enable SSL/TLS for database connections in production
- Consider rate limiting for API endpoints
- Implement proper logging and monitoring

## License

MIT

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
