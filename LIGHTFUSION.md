# Sun Ready - LightFusion API Integration

## Overview

Sun Ready integrates with the LightFusion API (running in Docker) to provide 3D solar project modeling and energy calculations. The LightFusion API handles:

- **3D Model Generation** from Google Earth data
- **Energy Requirement Calculations** (annual production)
- **Cost Estimations** for solar installations
- **Roof Plate Detection** and solar panel placement

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
