package cities

import (
	"backend/internal/models"
	"backend/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const indiaPostAPIURL = "https://dac.indiapost.gov.in/mypincode/home"
const indiaPostNextAction = "7f86f11580f1032603e2574b88b30708b68ae17e9d"

type indiaPostResponse struct {
	Data    []indiaPostOffice `json:"data"`
	Success bool              `json:"success"`
}

type indiaPostOffice struct {
	OfficeName    string `json:"office_name"`
	Pincode       string `json:"pincode"`
	DivisionName  string `json:"division_name"`
	RegionName    string `json:"region_name"`
	CircleName    string `json:"circle_name"`
	Taluk         string `json:"taluk"`
	DistrictName  string `json:"district_name"`
	StateName     string `json:"state_name"`
	WorkingHours  string `json:"working_hours"`
	ContactNumber string `json:"contact_number"`
	OfficeType    string `json:"office_type"`
	DeliveryStatus string `json:"delivery_status"`
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type BengaluruProvider struct {
	logger     *slog.Logger
	geoManager utils.GeoJsonManagerInterface
	httpClient httpClient
}

func NewBangaloreProvider(geoManager utils.GeoJsonManagerInterface, logger *slog.Logger) *BengaluruProvider {
	return &BengaluruProvider{
		logger:     logger,
		geoManager: geoManager,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
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

	if entity := p.getAssemblyConstituencyEntity(ctx, lat, lng); entity != nil {
		entities = append(entities, *entity)
	}

	if entity := p.getParliamentaryConstituencyEntity(ctx, lat, lng); entity != nil {
		entities = append(entities, *entity)
	}

	if entity := p.getPostOfficeEntity(ctx, lat, lng); entity != nil {
		entities = append(entities, *entity)
	}

	return entities, nil
}

func (p *BengaluruProvider) getGBAEntity(ctx context.Context, lat, lng float64) *models.Entity {
	attributes := append(
		utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "gba", nil, p.logger),
		utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "gba_ward", nil, p.logger)...,
	)

	entity := utils.BuildEntity(ctx, "GBA Corporation",
		"This information is unavailable for this address. This could be because the area is outside GBA Corporation limits.",
		nil,
		attributes, p.logger)

	return &entity
}

func (p *BengaluruProvider) getBBMPEntity(ctx context.Context, lat, lng float64) *models.Entity {
	attributes := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "bbmp", func(props map[string]interface{}) ([]models.Attribute, error) {
		wardNumber := getStringAttribute(props, "Ward Number", "ward_number", "")
		ward := getStringAttribute(props, "Ward Name", "ward_name", "")
		zone := getStringAttribute(props, "Zone", "zone", "")
		division := getStringAttribute(props, "Division", "division", "")
		subdivision := getStringAttribute(props, "Subdivision", "subdivision", "")

		return []models.Attribute{*wardNumber, *ward, *zone, *division, *subdivision}, nil
	}, p.logger)

	disclaimer := "This information is based on the 198-ward classification, which BBMP used as a reference until the GBA was formed."

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
		sectionAttributes := []models.Attribute{*section}

		if onmData, ok := props["onm"].(map[string]interface{}); ok {
			onmOffice := getAddressAttribute(onmData, "O&M Office", "om_office_name", "", "")
			sectionAttributes = append(sectionAttributes, *onmOffice)
		}

		return sectionAttributes, nil
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

	if dataValue, ok := data[key]; ok && dataValue != nil {
		value := fmt.Sprintf("%v", dataValue)
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

func (p *BengaluruProvider) getAssemblyConstituencyEntity(ctx context.Context, lat, lng float64) *models.Entity {
	attributes := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "assembly_constituency", nil, p.logger)

	entity := utils.BuildEntity(ctx, "Assembly Constituency",
		"This information is unavailable for this address. This could be because the area is outside Bengaluru Assembly Constituency limits.",
		nil,
		attributes, p.logger)

	return &entity
}

func (p *BengaluruProvider) getParliamentaryConstituencyEntity(ctx context.Context, lat, lng float64) *models.Entity {
	attributes := utils.ExtractAttributes(ctx, p.geoManager, lat, lng, p.Name(), "parliamentary_constituency", nil, p.logger)

	entity := utils.BuildEntity(ctx, "Parliamentary Constituency",
		"This information is unavailable for this address. This could be because the area is outside Bengaluru Parliamentary Constituency limits.",
		nil,
		attributes, p.logger)

	return &entity
}

func (p *BengaluruProvider) getPostOfficeEntity(ctx context.Context, lat, lng float64) *models.Entity {
	offices, err := p.fetchIndiaPostOffices(ctx, lat, lng)
	if err != nil {
		p.logger.WarnContext(ctx, "failed to fetch post office data",
			slog.Any("error", err),
			slog.Float64("lat", lat),
			slog.Float64("lng", lng),
		)
		entity := models.NewUnavailableEntity("Post Office",
			"This information is unavailable for this address. This could be because the India Post API is temporarily unavailable.",
			nil,
		)
		return &entity
	}

	if len(offices) == 0 {
		p.logger.InfoContext(ctx, "no post offices found for location",
			slog.Float64("lat", lat),
			slog.Float64("lng", lng),
		)
		entity := models.NewUnavailableEntity("Post Office",
			"No post office information found for this location.",
			nil,
		)
		return &entity
	}

	var attributes []models.Attribute
	for i, office := range offices {
		attributes = append(attributes, extractPostOfficeAttributes(office, i)...)
	}

	entity := utils.BuildEntity(ctx, "Post Office",
		"This information is unavailable for this address. This could be because the India Post API is temporarily unavailable.",
		nil,
		attributes, p.logger)

	return &entity
}

func (p *BengaluruProvider) fetchIndiaPostOffices(ctx context.Context, lat, lng float64) ([]indiaPostOffice, error) {
	payload := fmt.Sprintf("[[%f,%f]]", lat, lng)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, indiaPostAPIURL, bytes.NewBufferString(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Next-Action", indiaPostNextAction)
	req.Header.Set("content-type", "text/plain")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	offices, err := parseNextJSStreamResponse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return offices, nil
}

func parseNextJSStreamResponse(body []byte) ([]indiaPostOffice, error) {
	lines := strings.Split(string(body), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "1:") {
			jsonStr := strings.TrimPrefix(line, "1:")

			var response indiaPostResponse
			if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
				return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
			}

			if !response.Success {
				return nil, fmt.Errorf("API returned success=false")
			}

			return response.Data, nil
		}
	}

	return nil, fmt.Errorf("no data line found in Next.js stream response")
}

func extractPostOfficeAttributes(office indiaPostOffice, index int) []models.Attribute {
	var attributes []models.Attribute

	prefix := ""
	if index > 0 {
		prefix = fmt.Sprintf(" (Office %d)", index+1)
	}

	if office.OfficeName != "" {
		attributes = append(attributes, models.Attribute{
			Name:    fmt.Sprintf("Office Name%s", prefix),
			Value:   office.OfficeName,
			IsFound: true,
		})
	}

	if office.Pincode != "" {
		attributes = append(attributes, models.Attribute{
			Name:    fmt.Sprintf("PIN Code%s", prefix),
			Value:   office.Pincode,
			IsFound: true,
		})
	}

	if office.DivisionName != "" {
		attributes = append(attributes, models.Attribute{
			Name:    fmt.Sprintf("Division%s", prefix),
			Value:   office.DivisionName,
			IsFound: true,
		})
	}

	if office.RegionName != "" {
		attributes = append(attributes, models.Attribute{
			Name:    fmt.Sprintf("Region%s", prefix),
			Value:   office.RegionName,
			IsFound: true,
		})
	}

	if office.CircleName != "" {
		attributes = append(attributes, models.Attribute{
			Name:    fmt.Sprintf("Circle%s", prefix),
			Value:   office.CircleName,
			IsFound: true,
		})
	}

	if office.Taluk != "" {
		attributes = append(attributes, models.Attribute{
			Name:    fmt.Sprintf("Taluk%s", prefix),
			Value:   office.Taluk,
			IsFound: true,
		})
	}

	if office.DistrictName != "" {
		attributes = append(attributes, models.Attribute{
			Name:    fmt.Sprintf("District%s", prefix),
			Value:   office.DistrictName,
			IsFound: true,
		})
	}

	if office.StateName != "" {
		attributes = append(attributes, models.Attribute{
			Name:    fmt.Sprintf("State%s", prefix),
			Value:   office.StateName,
			IsFound: true,
		})
	}

	if office.OfficeType != "" {
		attributes = append(attributes, models.Attribute{
			Name:    fmt.Sprintf("Office Type%s", prefix),
			Value:   office.OfficeType,
			IsFound: true,
		})
	}

	if office.DeliveryStatus != "" {
		attributes = append(attributes, models.Attribute{
			Name:    fmt.Sprintf("Delivery Status%s", prefix),
			Value:   office.DeliveryStatus,
			IsFound: true,
		})
	}

	if office.WorkingHours != "" {
		attributes = append(attributes, models.Attribute{
			Name:    fmt.Sprintf("Working Hours%s", prefix),
			Value:   office.WorkingHours,
			IsFound: true,
		})
	}

	if office.ContactNumber != "" {
		attributes = append(attributes, models.Attribute{
			Name:    fmt.Sprintf("Contact Number%s", prefix),
			Value:   office.ContactNumber,
			IsFound: true,
		})
	}

	return attributes
}
