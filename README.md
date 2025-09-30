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

## API Endpoints

### Public Endpoints

#### Health Check
```bash
GET /health
```

#### Register
```bash
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "company_id": 1
}
```

#### Login
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### Protected Endpoints (Require JWT Token)

All protected endpoints require an `Authorization` header:
```
Authorization: Bearer <your-jwt-token>
```

#### Users

- `GET /api/users/{id}` - Get user by ID
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user
- `GET /api/users?company_id=1&limit=20&offset=0` - List users

#### Projects

- `POST /api/projects` - Create project
- `GET /api/projects/{id}` - Get project by ID
- `PUT /api/projects/{id}` - Update project
- `DELETE /api/projects/{id}` - Delete project
- `GET /api/projects?company_id=1&limit=20&offset=0` - List projects by company
- `GET /api/projects/user?user_id=1&limit=20&offset=0` - List projects by user

## Example Usage

### 1. Register a new user
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "company_id": 1
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "company_id": 1
  }
}
```

### 3. Create a project (with token)
```bash
curl -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "company_id": 1,
    "user_id": 1,
    "name": "Solar Installation Project",
    "description": "Residential solar panel installation",
    "status": "draft",
    "address": "123 Main St, City, State"
  }'
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `JWT_SECRET` | Secret key for JWT signing | Required |
| `PORT` | Server port | `8080` |
| `ENV` | Environment (development/production) | `development` |

## Database Schema

### Companies
- `id` - Primary key
- `name` - Company name
- `slug` - URL-friendly identifier
- `is_active` - Active status
- Timestamps: `created_at`, `updated_at`

### Users
- `id` - Primary key
- `email` - Unique email address
- `password` - Hashed password
- `first_name`, `last_name` - User name
- `company_id` - Foreign key to companies
- `type` - User type (admin, sales, client)
- `disabled` - Account status
- Timestamps: `created_at`, `updated_at`

### Projects
- `id` - Primary key
- `name` - Project name
- `description` - Project details
- `status` - Project status (draft, in_progress, completed, cancelled)
- `company_id` - Foreign key to companies
- `user_id` - Foreign key to users
- `address` - Project location
- Timestamps: `created_at`, `updated_at`

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
