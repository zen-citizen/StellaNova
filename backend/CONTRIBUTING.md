# Contributing a New City

This guide assumes that you already have the required boundaries and properties (as a GeoJSON) for that data that you want to add. If you need help finding this data, please look over [here](../README.md).

Adding a new city involves creating GeoJSON boundary files and implementing a city provider. We'll start with a simple example and then look at more complex cases.

> [!TIP]
> Look at the [Bengaluru provider](internal/cities/bengaluru.go) to see how transformers and city-specific utility functions are implemented

### Basic Approach: Using Pre-formatted Properties

If your GeoJSON already has well-formatted properties, you can use the default transformer (passing `nil`) and let the system handle everything automatically.

#### Step 1: Create GeoJSON with Properly Formatted Properties

Create the directory:
```
assets/geojson/stardew_valley/
```

Add `assets/geojson/stardew_valley/farm_district.geo.json`:
```json
{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "properties": {
        "District Name": "Pelican Town",
        "Mayor": "Lewis",
        "Population": "35 residents"
      },
      "geometry": {
        "type": "Polygon",
        "coordinates": [[
          [77.5, 12.9],
          [77.6, 12.9],
          [77.6, 13.0],
          [77.5, 13.0],
          [77.5, 12.9]
        ]]
      }
    }
  ]
}
```

The formatted keys are used as-is in the response.

#### Step 2: Create Simple City Provider

Create `internal/cities/stardew_valley.go`:

```go
package cities

import (
    "backend/internal/models"
    "backend/internal/utils"
    "context"
    "log/slog"
)

type StardewValleyProvider struct {
    logger     *slog.Logger
    geoManager *utils.GeoJsonManager
}

func NewStardewValleyProvider(geoManager *utils.GeoJsonManager, logger *slog.Logger) *StardewValleyProvider {
    return &StardewValleyProvider{
        logger:     logger,
        geoManager: geoManager,
    }
}

func (p *StardewValleyProvider) Name() string {
    return "stardew_valley"
}

func (p *StardewValleyProvider) FormattedName() string {
    return "Stardew Valley"
}

func (p *StardewValleyProvider) Bounds() *models.Bounds {
    return &models.Bounds{
        Northeast: models.Coordinate{Lat: 13.0, Lng: 77.7},
        Southwest: models.Coordinate{Lat: 12.9, Lng: 77.5},
    }
}

func (p *StardewValleyProvider) GetEntities(ctx context.Context, lat, lng float64) ([]models.Entity, error) {
    var entities []models.Entity
    
    if entity := p.getFarmDistrictEntity(ctx, lat, lng); entity != nil {
        entities = append(entities, *entity)
    }
    
    return entities, nil
}

func (p *StardewValleyProvider) getFarmDistrictEntity(ctx context.Context, lat, lng float64) *models.Entity {
    // Pass nil as transformer - properties will be mapped directly
    attributes := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "farm_district", nil, p.logger)
    
    entity := utils.BuildEntity(ctx, 
        "Farm District",
        "This location is outside Stardew Valley farm districts",
        nil,
        attributes, 
        p.logger)
    
    return &entity
}
```

#### Step 3: Register the City

Add to `internal/cities/registry.go`:
```go
r.cities["stardew_valley"] = NewStardewValleyProvider(geoManager, logger)
```

Now you query a point inside Pelican Town, you'll get:

```json
{
  "entities": [
    {
      "name": "Farm District",
      "is_available": true,
      "not_available_message": "This location is outside Stardew Valley farm districts",
      "attributes": [
        {
          "name": "District Name",
          "value": "Pelican Town"
        },
        {
          "name": "Mayor",
          "value": "Lewis"
        },
        {
          "name": "Population",
          "value": "35 residents"
        }
      ]
    }
  ]
}
```

#### How Empty Attributes and Entities Are Handled

The `IsFound` field in attributes is crucial for determining entity availability:

- When `ExtractAttributes` finds no matching polygon, it returns an empty slice
- When all attributes have `IsFound: false`, the entity is marked as unavailable
- `BuildEntity` automatically sets `is_available: false` when no attributes are found or all are not found
- The `not_available_message` is shown when `is_available` is false

For coordinates outside all boundaries:
```json
{
  "entities": [
    {
      "name": "Farm District",
      "is_available": false,
      "not_available_message": "This location is outside Stardew Valley farm districts",
      "attributes": []
    }
  ]
}
```

### Advanced: Using Custom Transformers

In reality, GeoJSON data is rarely perfectly formatted. You might get data from government sources with cryptic field names, or you might want to supplement with your own data. Having a transformer means you can handle messy data without constantly editing GeoJSON files.

Let's say your GeoJSON actually looks like this (more realistic):

```json
{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "properties": {
        "dst_nm": "Pelican Town",
        "admin": {
          "mayor": "Lewis",
          "term_start": "1990"
        },
        "stats": {
          "pop_2020": 30,
          "pop_2024": 35,
          "households": 12
        },
        "services": ["clinic", "general_store", "saloon"]
      },
      "geometry": {
        "type": "Polygon",
        "coordinates": [[...]]
      }
    }
  ]
}
```

This has abbreviated keys, nested objects, and raw data that needs formatting. Here's how to handle it with a transformer:

```go
func (p *StardewValleyProvider) getFarmDistrictEntity(ctx context.Context, lat, lng float64) *models.Entity {
    attributes := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "farm_district", 
        func(props map[string]interface{}) ([]models.Attribute, error) {
            var attrs []models.Attribute
            
            // Handle abbreviated field name
            if name, ok := props["dst_nm"].(string); ok {
                attrs = append(attrs, models.Attribute{
                    Name:    "District",
                    Value:   name,
                    IsFound: true,
                })
            }
            
            // Extract from nested object
            if admin, ok := props["admin"].(map[string]interface{}); ok {
                if mayor, ok := admin["mayor"].(string); ok {
                    if termStart, ok := admin["term_start"].(string); ok {
                        // Format the data nicely
                        attrs = append(attrs, models.Attribute{
                            Name:    "Administration",
                            Value:   fmt.Sprintf("Mayor %s (since %s)", mayor, termStart),
                            IsFound: true,
                        })
                    }
                }
            }
            
            // Process nested statistics and calculate current population
            if stats, ok := props["stats"].(map[string]interface{}); ok {
                if pop2024, ok := stats["pop_2024"].(float64); ok {
                    if households, ok := stats["households"].(float64); ok {
                        avgPerHousehold := pop2024 / households
                        attrs = append(attrs, models.Attribute{
                            Name:    "Population",
                            Value:   fmt.Sprintf("%d residents across %d households (%.1f per household)", 
                                    int(pop2024), int(households), avgPerHousehold),
                            IsFound: true,
                        })
                    }
                }
            }
            
            // Convert array to readable string
            if services, ok := props["services"].([]interface{}); ok {
                serviceList := make([]string, 0, len(services))
                for _, s := range services {
                    if svc, ok := s.(string); ok {
                        // Capitalize service names
                        serviceList = append(serviceList, strings.Title(svc))
                    }
                }
                if len(serviceList) > 0 {
                    attrs = append(attrs, models.Attribute{
                        Name:    "Available Services",
                        Value:   strings.Join(serviceList, ", "),
                        IsFound: true,
                    })
                }
            }
            
            return attrs, nil
        }, p.logger)
    
    entity := utils.BuildEntity(ctx, 
        "Farm District",
        "This location is outside Stardew Valley farm districts",
        nil,
        attributes, 
        p.logger)
    
    return &entity
}
```

Now the same messy GeoJSON produces clean output:

```json
{
  "entities": [
    {
      "name": "Farm District",
      "is_available": true,
      "not_available_message": "This location is outside Stardew Valley farm districts",
      "attributes": [
        {
          "name": "District",
          "value": "Pelican Town"
        },
        {
          "name": "Administration",
          "value": "Mayor Lewis (since 1990)"
        },
        {
          "name": "Population",
          "value": "35 residents across 12 households (2.9 per household)"
        },
        {
          "name": "Available Services",
          "value": "Clinic, General_store, Saloon"
        }
      ]
    }
  ]
}
```

### Things to look out for

1. **Set `IsFound` correctly**: Set to `true` when data exists, `false` when using defaults
2. **Handle type assertions safely**: Always check `ok` when casting interface{} values
3. **Keep transformers focused**: Each entity method should handle one logical boundary type
4. **Entity separation**: Each entity should represent one department/service (e.g., separate entities for electricity and water supply, not combined)


## Testing Your City Provider

Your city provider will be automatically tested when you run the test suite.

### How Testing Works

The [test suite](internal/cities/all_providers_test.go) automatically:
1. Discovers all registered city providers from the registry
2. Scans your GeoJSON directory to find all layers
3. Verifies your provider queries all layers
4. Tests your provider with various data scenarios
5. Ensures proper error handling and data quality

### Running Tests

```bash
# Test all city providers (including yours)
go test ./internal/cities -v

# Test only your city (for debugging)
go test ./internal/cities -v -run TestAllCityProviders/{your_city_name}
```

### What Gets Tested

The universal test suite verifies:

#### 1. Provider Properties
- Name matches directory name
- FormattedName is not empty
- Bounds are valid (NE > SW, within Earth's coordinates)

#### 2. Layer Coverage
- All `.geo.json` files in your directory are queried
- Provider doesn't skip any layers

#### 3. Data Scenarios (for each layer)
- **Nil response**: Point outside boundary
- **Valid data**: Properly formatted properties
- **Wrong types**: Strings instead of numbers, etc.
- **Missing fields**: Some properties absent
- **Null values**: Properties explicitly set to null
- **Nested issues**: Malformed nested structures

#### 4. Quality Checks
- `IsFound=true` → `Value` must not be empty
- `IsFound=false` → Should use default value (not empty)
- `IsAvailable=false` → `NotAvailableMessage` must not be empty
- No panics on any input
- All entities returned even when unavailable
