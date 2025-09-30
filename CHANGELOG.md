# Changelog

All notable changes to the SunReady API project will be documented in this file.

## [1.0.0] - 2025-09-30

### Added
- Initial project setup with Go 1.20
- Clean architecture with cmd/internal structure
- Repository pattern for data access layer
- Service pattern for business logic layer
- HTTP handlers for API endpoints
- GORM integration for database operations
- PostgreSQL database support
- JWT-based authentication
- User management APIs
- Company management APIs
- Project management APIs
- Docker support with Dockerfile
- Docker Compose for local development
- Database initialization scripts
- Comprehensive README documentation
- Makefile for common tasks
- Hot reload support with Air
- CORS middleware
- Health check endpoint
- Environment configuration with .env support

### Security
- Password hashing with bcrypt
- JWT token-based authentication
- Protected API endpoints with middleware

### Database Schema
- Companies table
- Users table with foreign key to companies
- Projects table with foreign keys to companies and users
- Proper indexes for performance

### API Endpoints
- POST /api/auth/register - User registration
- POST /api/auth/login - User login
- GET /api/users/{id} - Get user by ID
- PUT /api/users/{id} - Update user
- DELETE /api/users/{id} - Delete user
- GET /api/users - List users
- POST /api/projects - Create project
- GET /api/projects/{id} - Get project by ID
- PUT /api/projects/{id} - Update project
- DELETE /api/projects/{id} - Delete project
- GET /api/projects - List projects by company
- GET /api/projects/user - List projects by user
- GET /health - Health check endpoint
