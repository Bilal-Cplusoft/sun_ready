# SunReady API

A simplified solar project management API built with Go, featuring a clean architecture with repository and service patterns.

## Architecture

```
sunready/
├── cmd/
│   └── sunready/          # API Service
│       └── main.go
├── internal/
│   ├── database/          # Database connection
│   ├── client/            # Third party clients (Twilio, Lightfusion)
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


## Environment Variables
Inside env.example at /

## Security Notes

- Change `JWT_SECRET` in production
- Use strong passwords
- Enable SSL/TLS for database connections in production
- Consider rate limiting for API endpoints
- Implement proper logging and monitoring

## License

internal

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
