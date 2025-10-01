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

## [1.1.0] - 2025-10-01

### Added - Companies Module Enhancement
- Enhanced Company model with fields from API project:
  - Sales commission fields (min, max, default)
  - Baseline and baseline adder fields
  - Contract tag, referred by user ID
  - Credits, custom commissions, pricing mode
  - Validation and sanitization methods
- Company CRUD endpoints with Swagger documentation:
  - POST /api/companies - Create company
  - POST /api/companies/add - Add company with user migration
  - GET /api/companies - List companies (with pagination and filtering)
  - GET /api/companies/all - Get all companies
  - GET /api/companies/{id} - Get company by ID
  - GET /api/companies/slug/{slug} - Get company by slug
  - PUT /api/companies/{id} - Update company
  - DELETE /api/companies/{id} - Delete company
- Query parameter filtering:
  - Filter by referred_by_user_id
  - Pagination with limit and offset
- User repository enhancements:
  - FindByIDs() - Fetch multiple users
  - GetDescendantIDs() - Recursive CTE for user hierarchy
  - UpdateCompanyID() - Update user's company
- AddCompany endpoint features:
  - Creates company and migrates main user
  - Recursively migrates all descendant users
  - Sets up referral relationships
- Swagger/OpenAPI documentation:
  - Full API documentation at /swagger/index.html
  - Request/response schemas with examples
  - Authentication requirements documented
- Database migration script for new company fields

### Changed
- CompanyHandler now accepts UserService for user migration
- Updated .air.toml to auto-regenerate Swagger docs
- Enhanced docker-compose.yaml volume configuration

### Database
- Added migration: db/migrations/001_add_company_fields.sql
- Updated init.sql with new company schema fields
- Added index for referred_by_user_id

### Documentation
- Added Swagger annotations to all endpoints
- Auto-generated API documentation

## [1.2.0] - 2025-10-01

### Added - Deals Module
- Complete deals management system replicating API project functionality
- Deal model with comprehensive fields:
  - Lead, Project, System, and Hardware references
  - Panel and Inverter tracking
  - Financial details (costs, commissions, profit)
  - Status tracking (pending, approved, installed)
  - Consumption and production KWH
- Deal CRUD endpoints with Swagger documentation:
  - POST /api/deals - Create deal
  - GET /api/deals - List all deals (paginated)
  - GET /api/deals/{id} - Get deal by ID
  - GET /api/deals/uuid/{uuid} - Get deal by UUID
  - PUT /api/deals/{id} - Update deal
  - DELETE /api/deals/{id} - Delete deal
  - GET /api/deals/company/{company_id} - List deals by company
  - GET /api/deals/company/{company_id}/signed - List signed deals
  - POST /api/deals/{id}/archive - Archive deal
  - POST /api/deals/{id}/unarchive - Unarchive deal
- Deal repository with specialized queries:
  - List by company, sales, homeowner, project
  - Filter signed deals
  - Archive/unarchive functionality
- Deal service layer with validation
- Database migration for deals table with proper indexes

### Database
- Added migration: db/migrations/002_create_deals_table.sql
- Deals table with foreign keys to:
  - Projects (required)
  - Users (sales and homeowner)
  - Companies (required)
  - Leads, Systems, Hardware (optional)
- Indexes for performance on all foreign keys and status fields

### Files Created
- internal/models/deal.go - Deal model
- internal/repo/deal_repo.go - Deal repository
- internal/service/deal_service.go - Deal service
- internal/handler/deal_handler.go - Deal HTTP handlers
- db/migrations/002_create_deals_table.sql - Database schema

### Notes
- Deals support both lead-based (API project style) and project-based workflows
- Panel and inverter IDs tracked for hardware specifications
- Full financial tracking including target EPC, costs, and profit margins
- Status workflow: pending → approved → installed
- Archive functionality for soft deletion

## [1.3.0] - 2025-10-01

### Added - Leads Module (Simplified CRUD)
- Complete leads management system with simplified CRUD operations
- Lead model with comprehensive fields:
  - Location tracking (latitude, longitude, address)
  - Energy details (kWh usage, system size, panel count)
  - System components (panels, inverters, batteries)
  - Utility and tariff tracking
  - Roof and surface details
  - Production estimates
  - Workflow states (welcome call, financing, design, installation, PTO, etc.)
  - Financial details (electricity costs, incentives)
- Lead CRUD endpoints with Swagger documentation:
  - POST /api/leads - Create lead
  - GET /api/leads - List all leads (paginated)
  - GET /api/leads/{id} - Get lead by ID
  - PUT /api/leads/{id} - Update lead
  - PUT /api/leads/{id}/state - Update lead state
  - DELETE /api/leads/{id} - Delete lead
  - GET /api/leads/company/{company_id} - List leads by company
- Lead repository with specialized queries:
  - List by company, creator, state
  - State management
- Lead service layer with validation
- Database migration for leads table with proper indexes

### Database
- Added migration: db/migrations/003_create_leads_table.sql
- Leads table with foreign keys to:
  - Companies (required)
  - Users (creator, required)
  - Panels, Inverters, Utilities (optional)
- Indexes for performance on company, creator, state, and location

### Files Created
- internal/models/lead.go - Lead model with workflow states
- internal/repo/lead_repo.go - Lead repository
- internal/service/lead_service.go - Lead service
- internal/handler/lead_handler.go - Lead HTTP handlers
- db/migrations/003_create_leads_table.sql - Database schema

### Notes
- Simplified implementation focusing on CRUD operations
- No 3D graphics or Google Cloud integrations
- No scene image processing
- Workflow states tracked for full solar installation pipeline
- Supports lead states: Progress, Done, Errored, Initialized
- Lead sources: Legacy, Drone, Earth, Flyover

## [1.4.0] - 2025-10-01

### Added - Hybrid Lead Management with External API Integration
- **LightFUSION API Client** for external lead management
  - HTTP client for communicating with external API
  - Support for create, read, update, and list operations
  - Automatic authentication with API key
  - Error handling and response parsing
- **Hybrid Lead Service** with dual-mode operation
  - Local-only mode: All leads stored in local database
  - Hybrid mode: Leads synced with external LightFUSION API
  - Automatic sync status tracking (pending, synced, failed)
  - Bi-directional sync: local → external and external → local
- **Lead Model Enhancements**
  - `external_lead_id` - Reference to external API lead
  - `sync_status` - Track sync state (pending/synced/failed)
  - `last_synced_at` - Timestamp of last successful sync
  - Helper methods: `MarkSynced()`, `MarkSyncFailed()`
- **Environment Configuration**
  - `LIGHTFUSION_API` - External API base URL (optional)
  - `LIGHTFUSION_API_KEY` - API authentication key (optional)
  - Automatic detection: if URL is empty, uses local-only mode

### Changed
- Lead service now accepts LightFusionClient and useExternalAPI flag
- Lead creation/update automatically syncs with external API when enabled
- Lead listing fetches from external API and syncs to local database
- Updated database migration to include sync tracking fields

### How It Works
1. **Local-Only Mode** (default)
   - Set `LIGHTFUSION_API=` (empty) in .env
   - All leads stored and managed locally
   - No external API calls

2. **Hybrid Mode** (with external API)
   - Set `LIGHTFUSION_API=https://your-api.com` in .env
   - Set `LIGHTFUSION_API_KEY=your-key` in .env
   - Leads created locally are pushed to external API
   - External leads are pulled and synced to local database
   - Local database maintains reference to external lead ID
   - Sync status tracked for each lead

### Benefits
- **Flexibility**: Use existing LightFUSION API for complex operations
- **Local Control**: Maintain your own lead data for reporting/analytics
- **Resilience**: Local database continues working if external API is down
- **Sync Tracking**: Know which leads are synced vs. pending
- **No Code Changes**: Switch modes via environment variables

### Files Created/Modified
- `internal/client/lightfusion_client.go` - External API client
- `internal/service/lead_service.go` - Updated with hybrid logic
- `internal/models/lead.go` - Added sync tracking fields
- `db/migrations/003_create_leads_table.sql` - Updated schema
- `.env.example` - Added LightFUSION configuration
- `docker-compose.yaml` - Added environment variables

### Example Configuration

**Local-only mode:**
```bash
LIGHTFUSION_API=
LIGHTFUSION_API_KEY=
```

**Hybrid mode with external API:**
```bash
LIGHTFUSION_API=https://api.lightfusion.com
LIGHTFUSION_API_KEY=sk_live_abc123xyz
```
