# 3D Model and Lead Integration Documentation

## Overview

This document describes the integration between our internal lead management system and the LightFusion 3D model generation API. The integration allows leads to be associated with 3D solar project models created via the LightFusion API.

## Architecture

### Database Schema Changes

We've added the following fields to the `leads` table to track 3D model associations:

- `lightfusion_3d_project_id` - References the LightFusion API project ID
- `lightfusion_3d_house_id` - References the LightFusion API house ID  
- `model_3d_status` - Tracks the status of 3D model generation (pending, processing, completed, failed)
- `model_3d_created_at` - Timestamp when 3D model creation was initiated
- `model_3d_completed_at` - Timestamp when 3D model was completed

### Model Updates

The `Lead` model has been enhanced with methods to manage 3D model associations:

- `Has3DModel()` - Checks if a lead has an associated 3D model
- `SetLightFusion3DProject(projectID, houseID)` - Links a lead to a LightFusion 3D project
- `Update3DModelStatus(status)` - Updates the 3D model generation status
- `Is3DModelReady()` - Checks if the 3D model is ready for viewing

## API Endpoints

### Lead Management

#### Create Lead
```
POST /api/leads
```
Creates a new lead with optional 3D model generation.

**Request Body:**
```json
{
  "company_id": 1,
  "creator_id": 1,
  "latitude": 37.7749,
  "longitude": -122.4194,
  "address": "123 Solar St, San Francisco, CA 94102",
  "source": 2,
  "kwh_usage": 12000,
  "system_size": 10.5,
  "panel_count": 30,
  "create_3d_model": false
}
```

#### Get Lead
```
GET /api/leads/{id}
```
Retrieves a lead by ID, including 3D model status.

#### List Leads
```
GET /api/leads?has_3d_model=true&company_id=1&limit=20&offset=0
```
Lists leads with optional filters for 3D models.

#### Update Lead
```
PUT /api/leads/{id}
```
Updates lead information.

#### Delete Lead
```
DELETE /api/leads/{id}
```
Deletes a lead.

#### Sync 3D Model Status
```
POST /api/leads/{id}/sync-3d-status
```
Synchronizes the 3D model status with the LightFusion API.

### 3D Project Management

#### Create 3D Project
```
POST /api/projects/3d
```
Creates a 3D project via LightFusion and optionally links it to a lead.

**Request Body:**
```json
{
  "lead_id": 123,  // Optional: Link to existing lead
  "latitude": 37.7749,
  "longitude": -122.4194,
  "address": {
    "street": "123 Solar Street",
    "city": "San Francisco",
    "state": "CA",
    "Postalcode": "94102",
    "country": "USA"
  },
  "homeowner": {
    "email": "homeowner@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890"
  },
  "hardware": {
    "panel_id": 1,
    "inverter_id": 1,
    "storage_id": 1,
    "storage_quantity": 2
  },
  "consumption": [800, 850, 900, 950, 1000, 1050, 1100, 1150, 1200, 1250, 1300, 1350],
  "lse_id": 1,
  "period": "month",
  "target_solar_offset": 100,
  "mode": "max",
  "unit": "kwh",
  "company_id": 1,  // Required if creating new lead
  "creator_id": 1   // Required if creating new lead
}
```

**Response:**
```json
{
  "id": 123,
  "lead_id": 456,
  "status": "processing",
  "annual_production": 15000,
  "system_size": 10.5,
  "estimated_cost": 25000,
  "annual_savings": 2500,
  "message": "3D project created successfully. Processing in background."
}
```

#### Get Project Status
```
GET /api/projects/3d/{id}?house_id={house_id}
```
Retrieves the status of a 3D project and updates the associated lead if found.

#### Get Project Files
```
GET /api/projects/3d/{id}/files
```
Retrieves the 3D mesh files (JPG, OBJ, PLY, MTL) for a project.

## Integration Flow

### Creating a 3D Model for a Lead

1. **Option 1: Create Lead First**
   - Create a lead via `POST /api/leads`
   - Create a 3D project via `POST /api/projects/3d` with the `lead_id`
   - The lead is automatically updated with LightFusion project references

2. **Option 2: Create 3D Project Directly**
   - Create a 3D project via `POST /api/projects/3d` without a `lead_id`
   - A new lead is automatically created with the project information
   - The lead is linked to the LightFusion project

### Status Synchronization

The system automatically updates lead information when:
- A 3D project is created
- Project status is retrieved via `GET /api/projects/3d/{id}`
- Manual sync is triggered via `POST /api/leads/{id}/sync-3d-status`

Status values from LightFusion are mapped as follows:
- State 0 (Progress) → "processing"
- State 1 (Done) → "completed"
- State 2 (Errored) → "failed"
- State 3 (Initialized) → "processing"

### Data Synchronization

When retrieving project status, the following lead fields are automatically updated:
- `model_3d_status` - Current status of the 3D model
- `system_size` - System size from the 3D analysis
- `panel_count` - Number of panels from the 3D analysis
- `annual_production` - Annual energy production estimate

## Database Migrations

Apply the migration to add 3D model fields to the leads table:

```sql
-- File: db/migrations/001_add_3d_fields_to_leads.sql
ALTER TABLE leads 
ADD COLUMN IF NOT EXISTS lightfusion_3d_project_id INTEGER,
ADD COLUMN IF NOT EXISTS lightfusion_3d_house_id INTEGER,
ADD COLUMN IF NOT EXISTS model_3d_status VARCHAR(50),
ADD COLUMN IF NOT EXISTS model_3d_created_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS model_3d_completed_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_leads_lightfusion_3d_project_id ON leads(lightfusion_3d_project_id);
CREATE INDEX IF NOT EXISTS idx_leads_lightfusion_3d_house_id ON leads(lightfusion_3d_house_id);
CREATE INDEX IF NOT EXISTS idx_leads_model_3d_status ON leads(model_3d_status);
```

## Configuration

Ensure the following environment variables are set:

```env
# LightFusion API Configuration
LIGHTFUSION_API=http://localhost:8085
LIGHTFUSION_API_KEY=your_api_key_here
LIGHTFUSION_EMAIL=your_email@example.com
LIGHTFUSION_PASSWORD=your_password_here
```

## Error Handling

The integration includes robust error handling:
- Invalid lead IDs return 404 Not Found
- Failed LightFusion API calls are logged but don't break the lead creation flow
- Status synchronization failures are non-fatal and logged for debugging

## Best Practices

1. **Async Processing**: Consider implementing async workers for 3D model generation as it can take 30-120 seconds
2. **Status Polling**: Implement periodic status checks for "processing" models
3. **Error Recovery**: Failed 3D model generations should be retryable
4. **Caching**: Consider caching 3D model files locally to reduce API calls
5. **Monitoring**: Track 3D model generation success rates and processing times

## Example Use Cases

### Create a Lead with 3D Model
```bash
# 1. Create a lead
curl -X POST http://localhost:8080/api/leads \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": 1,
    "creator_id": 1,
    "latitude": 37.7749,
    "longitude": -122.4194,
    "address": "123 Solar St, San Francisco, CA 94102",
    "kwh_usage": 12000
  }'

# 2. Create 3D project for the lead (assuming lead_id is 1)
curl -X POST http://localhost:8080/api/projects/3d \
  -H "Content-Type: application/json" \
  -d '{
    "lead_id": 1,
    "latitude": 37.7749,
    "longitude": -122.4194,
    "address": {...},
    "homeowner": {...},
    "hardware": {...},
    "consumption": [...]
  }'

# 3. Check status
curl http://localhost:8080/api/projects/3d/123?house_id=456

# 4. Sync lead status
curl -X POST http://localhost:8080/api/leads/1/sync-3d-status
```

### List Leads with 3D Models
```bash
curl "http://localhost:8080/api/leads?has_3d_model=true&limit=10"
```

## Troubleshooting

### Common Issues

1. **3D Model Not Linking to Lead**
   - Ensure `lead_id` is provided when creating 3D project
   - Check that the lead exists before creating the project

2. **Status Not Updating**
   - Verify LightFusion API credentials are correct
   - Check network connectivity to LightFusion API
   - Review logs for API error responses

3. **Missing 3D Files**
   - Ensure the 3D model status is "completed" before accessing files
   - Check LightFusion API response for file URLs
   - Verify file download permissions

## Future Enhancements

1. **Webhook Support**: Implement webhooks for real-time status updates from LightFusion
2. **Batch Processing**: Add support for creating multiple 3D models in bulk
3. **File Storage**: Implement local storage and CDN integration for 3D files
4. **Analytics**: Add tracking for model generation metrics and success rates
5. **Notifications**: Send notifications when 3D models are ready