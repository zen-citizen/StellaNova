package cities

import (
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"fmt"
	"log/slog"
)

type BengaluruProvider struct {
	logger     *slog.Logger
	geoManager *utils.GeoJsonManager
}

func NewBangaloreProvider(geoManager *utils.GeoJsonManager, logger *slog.Logger) *BengaluruProvider {
	return &BengaluruProvider{
		logger:     logger,
		geoManager: geoManager,
	}
}

func (p *BengaluruProvider) FormattedName() string {
	return "Bengaluru"
}

func (p *BengaluruProvider) Name() string {
	return "bengaluru"
}

func (p *BengaluruProvider) Bounds() *models.Bounds {
	return &models.Bounds{
		Northeast: models.Coordinate{Lat: 13.2, Lng: 77.8},
		Southwest: models.Coordinate{Lat: 12.8, Lng: 77.4},
	}
}

func (p *BengaluruProvider) GetEntities(ctx context.Context, lat, lng float64) ([]models.Entity, error) {
	p.logger.InfoContext(ctx, "getting entities", slog.String("city", p.Name()))

	var entities []models.Entity

	if entity := p.getGBAEntity(ctx, lat, lng); entity != nil {
		entities = append(entities, *entity)
	}

	if entity := p.getBBMPEntity(ctx, lat, lng); entity != nil {
		entities = append(entities, *entity)
	}

	if entity := p.getBDAEntity(ctx, lat, lng); entity != nil {
		entities = append(entities, *entity)
	}

	if entity := p.getRevenueClassificationEntity(ctx, lat, lng); entity != nil {
		entities = append(entities, *entity)
	}

	if entity := p.getRevenueOfficesEntity(ctx, lat, lng); entity != nil {
		entities = append(entities, *entity)
	}

	if entity := p.getBESCOMEntity(ctx, lat, lng); entity != nil {
		entities = append(entities, *entity)
	}

	if entity := p.getBWSSBEntity(ctx, lat, lng); entity != nil {
		entities = append(entities, *entity)
	}

	if entity := p.getPoliceEntity(ctx, lat, lng); entity != nil {
		entities = append(entities, *entity)
	}

	return entities, nil
}

func (p *BengaluruProvider) getGBAEntity(ctx context.Context, lat, lng float64) *models.Entity {
	attributes := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "gba", func(props map[string]interface{}) ([]models.Attribute, error) {
		attribute := getStringAttribute(props, "Corporation Name", "name", "")

		return []models.Attribute{*attribute}, nil
	}, p.logger)

	entity := utils.BuildEntity(ctx, "GBA Corporation",
		"This information is unavailable for this address. This could be because the area is outside GBA Corporation limits.",
		nil,
		attributes, p.logger)

	return &entity
}

func (p *BengaluruProvider) getBBMPEntity(ctx context.Context, lat, lng float64) *models.Entity {
	attributes := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "bbmp", func(props map[string]interface{}) ([]models.Attribute, error) {
		ward := getStringAttribute(props, "Ward Name", "ward_name", "")
		wardNumber := getStringAttribute(props, "Ward Number", "ward_number", "")
		zone := getStringAttribute(props, "Zone", "zone", "")
		division := getStringAttribute(props, "Division", "division", "")
		subdivision := getStringAttribute(props, "Subdivision", "subdivision", "")

		return []models.Attribute{*ward, *wardNumber, *zone, *division, *subdivision}, nil
	}, p.logger)

	disclaimer := "This information is based on the 198-ward classification, which BBMP still uses as a reference — even though it's no longer the official structure."

	entity := utils.BuildEntity(ctx, "BBMP Information",
		"This information is unavailable for this address. This could be because the area is outside BBMP limits.",
		&disclaimer,
		attributes, p.logger)

	return &entity
}

func (p *BengaluruProvider) getBDAEntity(ctx context.Context, lat, lng float64) *models.Entity {
	attributes := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "bda", func(props map[string]interface{}) ([]models.Attribute, error) {
		layoutName := getStringAttribute(props, "BDA Layout Name", "name", "")
		layoutNumber := getStringAttribute(props, "BDA Layout Number", "layout_number", "")

		return []models.Attribute{*layoutName, *layoutNumber}, nil
	}, p.logger)

	entity := utils.BuildEntity(ctx, "BDA Information",
		"This information is unavailable for this address. This could be because the area is outside BDA limits.",
		nil,
		attributes, p.logger)

	return &entity
}

func (p *BengaluruProvider) getRevenueClassificationEntity(ctx context.Context, lat, lng float64) *models.Entity {
	var attributes []models.Attribute

	districtAttrs := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "revenue_district", func(props map[string]interface{}) ([]models.Attribute, error) {
		district := getStringAttribute(props, "District", "name", "")
		return []models.Attribute{*district}, nil
	}, p.logger)
	attributes = append(attributes, districtAttrs...)

	talukAttrs := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "revenue_taluk", func(props map[string]interface{}) ([]models.Attribute, error) {
		taluk := getStringAttribute(props, "Taluk", "name", "")
		return []models.Attribute{*taluk}, nil
	}, p.logger)
	attributes = append(attributes, talukAttrs...)

	hobliAttrs := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "revenue_hobli", func(props map[string]interface{}) ([]models.Attribute, error) {
		hobli := getStringAttribute(props, "Hobli", "name", "")
		return []models.Attribute{*hobli}, nil
	}, p.logger)
	attributes = append(attributes, hobliAttrs...)

	villageAttrs := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "revenue_village", func(props map[string]interface{}) ([]models.Attribute, error) {
		village := getStringAttribute(props, "Village", "name", "")
		return []models.Attribute{*village}, nil
	}, p.logger)
	attributes = append(attributes, villageAttrs...)

	entity := utils.BuildEntity(ctx, "Revenue Classification",
		"This information is unavailable for this address. This could be because the area is outside Bengaluru Urban district.",
		nil,
		attributes, p.logger)

	return &entity
}

func (p *BengaluruProvider) getRevenueOfficesEntity(ctx context.Context, lat, lng float64) *models.Entity {
	attributes := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "revenue_office", func(props map[string]interface{}) ([]models.Attribute, error) {
		var sro *models.Attribute
		if sroData, ok := props["sro"].(map[string]interface{}); !ok {
			sro = &models.Attribute{
				Name:    "SRO",
				Value:   "Not available",
				Address: nil,
				IsFound: false,
			}
		} else {
			sro = getAddressAttribute(sroData, "SRO", "name", "", "")
		}

		var dro *models.Attribute
		if droData, ok := props["dro"].(map[string]interface{}); !ok {
			dro = &models.Attribute{
				Name:    "DRO",
				Value:   "Not available",
				Address: nil,
				IsFound: false,
			}
		} else {
			dro = getAddressAttribute(droData, "DRO", "name", "", "")
		}

		return []models.Attribute{*sro, *dro}, nil
	}, p.logger)

	entity := utils.BuildEntity(ctx, "Revenue Offices",
		"This information is unavailable for this address. This could be because the area is outside Bengaluru Urban district.",
		nil,
		attributes, p.logger)

	return &entity
}

func (p *BengaluruProvider) getBESCOMEntity(ctx context.Context, lat, lng float64) *models.Entity {
	var attributes []models.Attribute

	divisionAttrs := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "bescom_division", func(props map[string]interface{}) ([]models.Attribute, error) {
		division := getStringAttribute(props, "Division", "name", "")
		return []models.Attribute{*division}, nil
	}, p.logger)
	attributes = append(attributes, divisionAttrs...)

	subdivisionAttrs := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "bescom_subdivision", func(props map[string]interface{}) ([]models.Attribute, error) {
		subdivision := getStringAttribute(props, "Subdivision", "name", "")
		return []models.Attribute{*subdivision}, nil
	}, p.logger)
	attributes = append(attributes, subdivisionAttrs...)

	sectionAttrs := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "bescom_section", func(props map[string]interface{}) ([]models.Attribute, error) {
		section := getStringAttribute(props, "Section", "name", "")

		onmOffice := &models.Attribute{}
		if onmData, ok := props["onm"].(map[string]interface{}); ok {
			onmOffice = getAddressAttribute(onmData, "O&M Office", "om_office_name", "", "")
		}

		return []models.Attribute{*section, *onmOffice}, nil
	}, p.logger)
	attributes = append(attributes, sectionAttrs...)

	if len(attributes) == 0 {
		return nil
	}

	disclaimer := "O&M office data may need verification."
	entity := utils.BuildEntity(ctx, "Electricity (BESCOM)",
		"This information is unavailable for this address. This could be because the area is outside BESCOM limits.",
		&disclaimer,
		attributes, p.logger)

	return &entity
}

func (p *BengaluruProvider) getBWSSBEntity(ctx context.Context, lat, lng float64) *models.Entity {
	var attributes []models.Attribute

	divisionAttrs := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "bwssb_division", func(props map[string]interface{}) ([]models.Attribute, error) {
		division := getStringAttribute(props, "Division", "name", "")
		return []models.Attribute{*division}, nil
	}, p.logger)
	attributes = append(attributes, divisionAttrs...)

	subdivisionAttrs := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "bwssb_subdivision", func(props map[string]interface{}) ([]models.Attribute, error) {
		subdivision := getStringAttribute(props, "Subdivision", "name", "")
		return []models.Attribute{*subdivision}, nil
	}, p.logger)
	attributes = append(attributes, subdivisionAttrs...)

	serviceStationAttrs := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "bwssb_service_station", func(props map[string]interface{}) ([]models.Attribute, error) {
		serviceStation := getStringAttribute(props, "Service Station", "name", "")
		return []models.Attribute{*serviceStation}, nil
	}, p.logger)
	attributes = append(attributes, serviceStationAttrs...)

	entity := utils.BuildEntity(ctx, "Water Supply (BWSSB)",
		"This information is unavailable for this address. This could be because the area is outside BWSSB service limits.",
		nil,
		attributes, p.logger)

	return &entity
}

func (p *BengaluruProvider) getPoliceEntity(ctx context.Context, lat, lng float64) *models.Entity {
	var attributes []models.Attribute
	policeStationAttr := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "police_city", func(props map[string]interface{}) ([]models.Attribute, error) {
		policeStation := getAddressAttribute(props, "Police Station", "name", "", "")
		return []models.Attribute{*policeStation}, nil
	}, p.logger)
	attributes = append(attributes, policeStationAttr...)

	trafficStationAttr := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "police_traffic", func(props map[string]interface{}) ([]models.Attribute, error) {
		trafficStation := getAddressAttribute(props, "Traffic Police", "name", "", "")
		return []models.Attribute{*trafficStation}, nil
	}, p.logger)
	attributes = append(attributes, trafficStationAttr...)

	entity := utils.BuildEntity(ctx, "Police Jurisdiction",
		"This information is unavailable for this address. This could be because the area is outside of Bengaluru City and Traffic police limits.",
		nil,
		attributes, p.logger)

	return &entity
}

func getStringAttribute(data map[string]interface{}, name, key, valueDefault string) *models.Attribute {
	if valueDefault == "" {
		valueDefault = "Not available"
	}

	value := valueDefault
	if dataValue, ok := data[key]; ok && dataValue != nil {
		value = fmt.Sprintf("%v", dataValue)
		return &models.Attribute{
			Name:    name,
			Value:   value,
			IsFound: true,
		}
	} else {
		return &models.Attribute{
			Name:    name,
			Value:   valueDefault,
			IsFound: false,
		}
	}
}

func getAddressAttribute(data map[string]interface{}, name, key, valueDefault, addressDefault string) *models.Attribute {
	attribute := getStringAttribute(data, name, key, valueDefault)
	if !attribute.IsFound {
		return attribute
	}

	if addressDefault == "" {
		addressDefault = "Address not available"
	}

	address := models.Address{
		Text: addressDefault,
	}
	if places, ok := data["places"].([]interface{}); ok && places != nil && len(places) != 0 {
		if firstAddress, firstAddressOk := places[0].(map[string]interface{}); firstAddressOk {
			if text, exists := firstAddress["formattedAddress"]; exists {
				address.Text = fmt.Sprintf("%v", text)
				if link, linkExists := firstAddress["googleMapsUri"]; linkExists {
					linkStr := fmt.Sprintf("%v", link)
					address.Link = &linkStr
				}
			}
		}
	}

	attribute.Address = &address
	return attribute
}
