# StellaNova backend

## Overview

This service provides location-based administrative and utility information for addresses in various cities. Give it a latitude and longitude, and it tells you which administrative zones, utility service areas, and jurisdictions that point falls under.

The service works by checking your coordinates against various GeoJSON boundary files to figure out which polygons contain your location. Each city has its own set of boundaries representing different administrative or service areas.

> [!NOTE]
> Look [here](CONTRIBUTING.md) for the guide to contributing.

## Running Locally

Prerequisites: Go 1.23+ installed on your machine

```bash
# Clone the repo and navigate to backend
cd backend

# Install dependencies
go mod download

# Run the server
go run cmd/server/main.go
```

The server starts on port 8080 by default. You can change this with the PORT environment variable.


Alternatively, you can use the [Dockerfile](Dockerfile)

To test if it's working:
```bash
curl "http://localhost:8080/api/v1/entities?lat=12.970110&lng=77.644720&city=bengaluru"
```

## Architecture

### Request Flow

Here's how a request moves through the layers:

1. **Router** (`router.go`) - Routes the request to the appropriate handler
2. **Handler** (`handlers/entities.go`) - Validates query parameters (lat, lng, city)
3. **Service** (`services/entities.go`) - Orchestrates the business logic
4. **City Provider** (`cities/{city}.go`) - Implements city-specific logic
5. **Utils** (`utils/geojson.go`, `utils/entity_builder.go`) - Handles GeoJSON queries and response building

### Key Parts

#### GeoJSON Manager (`internal/utils/geojson.go`)

Loads all GeoJSON files into memory at startup and provides point-in-polygon queries:

- Loads from `assets/geojson/{city}/*.geo.json`
- Files must end with `.geo.json`
- Provides `QueryPoint(lat, lng, city, layer)` to check if a point falls within any polygon

#### Entity Builder (`internal/utils/entity_builder.go`)

Provides utilities to extract attributes from GeoJSON properties and build response entities:

- `ExtractAttributes()` - Queries GeoJSON and transforms properties into attributes
- `BuildEntity()` - Creates an Entity with proper availability status

#### City Provider Interface (`internal/cities/provider.go`)

Every city must implement this interface
