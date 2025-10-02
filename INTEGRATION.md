# Sun Ready - LightFusion API Integration

## Overview

Sun Ready integrates with the LightFusion API (running in Docker) to provide 3D solar project modeling and energy calculations. The LightFusion API handles:

- **3D Model Generation** from Google Earth data
- **Energy Requirement Calculations** (annual production)
- **Cost Estimations** for solar installations
- **Roof Plate Detection** and solar panel placement

## Architecture

```
┌─────────────────┐      HTTP Requests       ┌──────────────────────┐
│   Sun Ready     │ ───────────────────────> │  LightFusion API     │
│   (Port 8080)   │                           │  (Port 8085)         │
│                 │ <─────────────────────── │  Docker Container    │
└─────────────────┘      JSON Responses       └──────────────────────┘
                                                        │
                                                        ├─> Google Earth API
                                                        ├─> Three.js Rendering
                                                        ├─> Mesh Generation
                                                        └─> Energy Calculations
```

## Running Services

### LightFusion API (Already Running)
```bash
docker ps | grep api
# Should show:
# - api-api-1 (Port 8085)
# - api-api_meshearth-1
# - api-edge-1 (Port 8081)
# - rabbitmq
# - redis
```

### Sun Ready API
```bash
cd /home/saint/lighthouse/new/sun_ready
# Update .env file
cp .env.example .env
# Edit .env and set LIGHTFUSION_API=http://localhost:8085

# Run the server
go run cmd/sunready/main.go
# Or with Air for hot reload
air
```

## API Endpoints

### Create 3D Solar Project

**Endpoint**: `POST /api/projects/3d`

**Description**: Creates a 3D model from Google Earth data and calculates energy requirements and costs.

**Request Body**:
```json
{
  "latitude": 37.7749,
  "longitude": -122.4194,
  "address": {
    "street": "123 Solar Street",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94102",
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
  "unit": "kwh"
}
```

**Response** (201 Created):
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

### Get Project Status

**Endpoint**: `GET /api/projects/3d/{id}`

**Description**: Retrieves the status and details of a 3D solar project.

**Response** (200 OK):
```json
{
  "id": 123,
  "lead_id": 456,
  "status": "completed",
  "annual_production": 15000,
  "system_size": 10.5,
  "estimated_cost": 25000,
  "annual_savings": 2500
}
```

## How It Works

### 1. User Request
User sends location (latitude/longitude) and energy consumption data to Sun Ready API.

### 2. Sun Ready → LightFusion
Sun Ready forwards the request to LightFusion API running on `http://localhost:8085`.

### 3. LightFusion Processing
LightFusion API:
- Downloads 3D mesh from Google Earth using the coordinates
- Uses headless Chrome to navigate Google Earth
- Captures screenshots from multiple angles
- Detects roof plates for solar panel placement
- Calculates energy production based on:
  - Geographic location (sun exposure)
  - Roof orientation and tilt
  - Panel specifications
  - Historical weather data
- Estimates costs based on system size and hardware

### 4. Background Processing
The 3D model generation happens asynchronously:
- Initial response returns immediately with `status: "processing"`
- Use the status endpoint to check progress
- Status changes: `processing` → `completed` or `failed`

### 5. Response to User
Sun Ready returns the results including:
- Annual energy production (kWh)
- System size (kW)
- Estimated installation cost
- Annual savings

## Environment Variables

```bash
# Database
DATABASE_URL=postgres://sunready:sunready@localhost:5433/sunready?sslmode=disable

# JWT Authentication
JWT_SECRET=your-secret-key-change-in-production

# Server Port
PORT=8080

# LightFusion API Integration
LIGHTFUSION_API=http://localhost:8085
LIGHTFUSION_API_KEY=  # Optional, leave empty for local Docker
```

## Testing

### Test with cURL

```bash
# Create a 3D project
curl -X POST http://localhost:8080/api/projects/3d \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": 37.7749,
    "longitude": -122.4194,
    "address": {
      "street": "123 Solar St",
      "city": "San Francisco",
      "state": "CA",
      "postal_code": "94102",
      "country": "USA"
    },
    "homeowner": {
      "email": "test@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "phone": "+1234567890"
    },
    "hardware": {
      "panel_id": 1,
      "inverter_id": 1
    },
    "consumption": [1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000],
    "lse_id": 1,
    "period": "month",
    "target_solar_offset": 100,
    "unit": "kwh"
  }'

# Check project status (replace {id} with actual ID from response)
curl http://localhost:8080/api/projects/3d/{id}
```

## Troubleshooting

### LightFusion API Not Responding
```bash
# Check if container is running
docker ps | grep api-api-1

# Check container logs
docker logs api-api-1

# Restart if needed
cd /home/saint/lighthouse/new/api
docker-compose restart api
```

### Connection Refused
- Ensure LightFusion API is running on port 8085
- Check firewall settings
- Verify `LIGHTFUSION_API` environment variable is set correctly

### Timeout Errors
- 3D model generation can take 30-120 seconds
- The timeout is set to 120 seconds in the client
- Check LightFusion API logs for processing errors

## Next Steps

1. ✅ Integration with LightFusion API complete
2. ✅ 3D project creation endpoint implemented
3. ✅ Status checking endpoint implemented
4. 🔄 Add webhook support for async notifications
5. 🔄 Store project results in Sun Ready database
6. 🔄 Add authentication to 3D project endpoints
7. 🔄 Create frontend UI for project visualization
