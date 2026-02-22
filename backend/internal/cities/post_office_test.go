package cities

import (
	"backend/internal/models"
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPClient implements httpClient interface for testing
type mockHTTPClient struct {
	response    *http.Response
	err         error
	statusCode  int
	body        string
	requestBody string
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Capture request body for verification
	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		m.requestBody = string(body)
	}

	if m.err != nil {
		return nil, m.err
	}

	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(strings.NewReader(m.body)),
	}, nil
}

func TestParseNextJSStreamResponse(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		expectError bool
		expectCount int
	}{
		{
			name: "valid response with single office",
			body: `0:{"a":"$@1","f":"","b":"vE0mwr0TbQKy9EsanIL7Y"}
1:{"data":[{"office_name":"Viveknagar S.O","pincode":"560047","division_name":"Bengaluru East","region_name":"Bangalore HQ","circle_name":"Karnataka","taluk":"Bangalore","district_name":"Bengaluru","state_name":"Karnataka","working_hours":"09:00-17:00","contact_number":"080-25301234","office_type":"S.O","delivery_status":"Delivery"}],"success":true}`,
			expectError: false,
			expectCount: 1,
		},
		{
			name: "valid response with multiple offices",
			body: `0:{"a":"$@1","f":"","b":"vE0mwr0TbQKy9EsanIL7Y"}
1:{"data":[{"office_name":"Office 1","pincode":"560001"},{"office_name":"Office 2","pincode":"560002"}],"success":true}`,
			expectError: false,
			expectCount: 2,
		},
		{
			name: "response with success=false",
			body: `0:{"a":"$@1","f":"","b":"vE0mwr0TbQKy9EsanIL7Y"}
1:{"data":[],"success":false}`,
			expectError: true,
			expectCount: 0,
		},
		{
			name: "response with no data line",
			body: `0:{"a":"$@1","f":"","b":"vE0mwr0TbQKy9EsanIL7Y"}`,
			expectError: true,
			expectCount: 0,
		},
		{
			name: "response with malformed JSON",
			body: `1:{invalid json}`,
			expectError: true,
			expectCount: 0,
		},
		{
			name: "response with empty data array",
			body: `1:{"data":[],"success":true}`,
			expectError: false,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offices, err := parseNextJSStreamResponse([]byte(tt.body))

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(offices) != tt.expectCount {
					t.Errorf("expected %d offices, got %d", tt.expectCount, len(offices))
				}
			}
		})
	}
}

func TestExtractPostOfficeAttributes(t *testing.T) {
	tests := []struct {
		name            string
		office          indiaPostOffice
		index           int
		expectMinAttrs  int
		expectedNames   []string
		expectedValues  map[string]string
	}{
		{
			name: "fully populated office",
			office: indiaPostOffice{
				OfficeName:     "Viveknagar S.O",
				Pincode:        "560047",
				DivisionName:   "Bengaluru East",
				RegionName:     "Bangalore HQ",
				CircleName:     "Karnataka",
				Taluk:          "Bangalore",
				DistrictName:   "Bengaluru",
				StateName:      "Karnataka",
				WorkingHours:   "09:00-17:00",
				ContactNumber:  "080-25301234",
				OfficeType:     "S.O",
				DeliveryStatus: "Delivery",
			},
			index:          0,
			expectMinAttrs: 11,
			expectedNames:  []string{"Office Name", "PIN Code", "Division", "State"},
			expectedValues: map[string]string{
				"Office Name": "Viveknagar S.O",
				"PIN Code":    "560047",
				"Division":    "Bengaluru East",
			},
		},
		{
			name: "minimally populated office",
			office: indiaPostOffice{
				OfficeName: "Test Office",
				Pincode:    "560001",
			},
			index:          0,
			expectMinAttrs: 2,
			expectedNames:  []string{"Office Name", "PIN Code"},
			expectedValues: map[string]string{
				"Office Name": "Test Office",
				"PIN Code":    "560001",
			},
		},
		{
			name: "second office with prefix",
			office: indiaPostOffice{
				OfficeName: "Second Office",
				Pincode:    "560002",
			},
			index:          1,
			expectMinAttrs: 2,
			expectedNames:  []string{"Office Name (Office 2)", "PIN Code (Office 2)"},
			expectedValues: map[string]string{
				"Office Name (Office 2)": "Second Office",
				"PIN Code (Office 2)":    "560002",
			},
		},
		{
			name:            "empty office",
			office:          indiaPostOffice{},
			index:           0,
			expectMinAttrs:  0,
			expectedNames:   []string{},
			expectedValues:  map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := extractPostOfficeAttributes(tt.office, tt.index)

			if len(attrs) < tt.expectMinAttrs {
				t.Errorf("expected at least %d attributes, got %d", tt.expectMinAttrs, len(attrs))
			}

			// Check expected attribute names exist
			attrNames := make(map[string]bool)
			for _, attr := range attrs {
				attrNames[attr.Name] = true
				if !attr.IsFound {
					t.Errorf("attribute '%s' should have IsFound=true", attr.Name)
				}
			}

			for _, name := range tt.expectedNames {
				if !attrNames[name] {
					t.Errorf("expected attribute '%s' not found", name)
				}
			}

			// Check expected values
			for _, attr := range attrs {
				if expectedVal, ok := tt.expectedValues[attr.Name]; ok {
					if attr.Value != expectedVal {
						t.Errorf("attribute '%s': expected value '%s', got '%s'", attr.Name, expectedVal, attr.Value)
					}
				}
			}
		})
	}
}

func TestGetPostOfficeEntity(t *testing.T) {
	ctx := context.Background()
	logger := slog.Default()

	tests := []struct {
		name           string
		mockResponse   string
		mockStatus     int
		mockError      error
		expectAvail    bool
		expectAttrCnt  int
		expectEntityNm string
	}{
		{
			name: "successful single office",
			mockResponse: `0:{"a":"$@1","f":"","b":"test"}
1:{"data":[{"office_name":"Test S.O","pincode":"560047","division_name":"Bangalore East"}],"success":true}`,
			mockStatus:     http.StatusOK,
			expectAvail:    true,
			expectAttrCnt:  3,
			expectEntityNm: "Post Office",
		},
		{
			name: "successful multiple offices",
			mockResponse: `0:{"a":"$@1","f":"","b":"test"}
1:{"data":[{"office_name":"Office 1","pincode":"560001"},{"office_name":"Office 2","pincode":"560002"}],"success":true}`,
			mockStatus:     http.StatusOK,
			expectAvail:    true,
			expectAttrCnt:  4, // 2 attributes per office
			expectEntityNm: "Post Office",
		},
		{
			name:           "API returns error",
			mockResponse:   "",
			mockStatus:     http.StatusInternalServerError,
			expectAvail:    false,
			expectEntityNm: "Post Office",
		},
		{
			name:           "API returns empty data",
			mockResponse:   `1:{"data":[],"success":true}`,
			mockStatus:     http.StatusOK,
			expectAvail:    false,
			expectEntityNm: "Post Office",
		},
		{
			name:           "API returns success false",
			mockResponse:   `1:{"data":[],"success":false}`,
			mockStatus:     http.StatusOK,
			expectAvail:    false,
			expectEntityNm: "Post Office",
		},
		{
			name:           "malformed response",
			mockResponse:   `not a valid response`,
			mockStatus:     http.StatusOK,
			expectAvail:    false,
			expectEntityNm: "Post Office",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				statusCode: tt.mockStatus,
				body:       tt.mockResponse,
				err:        tt.mockError,
			}

			provider := &BengaluruProvider{
				logger:     logger,
				geoManager: NewInstrumentedMockGeoJsonManager(),
				httpClient: mockClient,
			}

			entity := provider.getPostOfficeEntity(ctx, 12.97, 77.59)

			if entity == nil {
				t.Fatal("expected entity, got nil")
			}

			if entity.Name != tt.expectEntityNm {
				t.Errorf("expected entity name '%s', got '%s'", tt.expectEntityNm, entity.Name)
			}

			if entity.IsAvailable != tt.expectAvail {
				t.Errorf("expected IsAvailable=%v, got %v", tt.expectAvail, entity.IsAvailable)
			}

			if tt.expectAvail && len(entity.Attributes) < tt.expectAttrCnt {
				t.Errorf("expected at least %d attributes, got %d", tt.expectAttrCnt, len(entity.Attributes))
			}

			if !tt.expectAvail && entity.NotAvailableMessage == "" {
				t.Error("unavailable entity should have NotAvailableMessage")
			}
		})
	}
}

func TestGetPostOfficeEntity_RequestFormat(t *testing.T) {
	ctx := context.Background()
	logger := slog.Default()

	mockClient := &mockHTTPClient{
		statusCode: http.StatusOK,
		body:       `1:{"data":[{"office_name":"Test"}],"success":true}`,
	}

	provider := &BengaluruProvider{
		logger:     logger,
		geoManager: NewInstrumentedMockGeoJsonManager(),
		httpClient: mockClient,
	}

	lat := 12.963492
	lng := 77.613344

	_ = provider.getPostOfficeEntity(ctx, lat, lng)

	// Verify request body format
	expectedPayload := "[[12.963492,77.613344]]"
	if mockClient.requestBody != expectedPayload {
		t.Errorf("expected request body '%s', got '%s'", expectedPayload, mockClient.requestBody)
	}
}

func TestFetchIndiaPostOffices_ContextCancellation(t *testing.T) {
	logger := slog.Default()

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Simulate context cancellation error
	mockClient := &mockHTTPClient{
		statusCode: http.StatusOK,
		body:       `1:{"data":[],"success":true}`,
		err:        context.Canceled,
	}

	provider := &BengaluruProvider{
		logger:     logger,
		geoManager: NewInstrumentedMockGeoJsonManager(),
		httpClient: mockClient,
	}

	_, err := provider.fetchIndiaPostOffices(ctx, 12.97, 77.59)
	if err == nil {
		t.Error("expected error due to cancelled context")
	}
}

func TestGetPostOfficeEntity_IntegrationWithGetEntities(t *testing.T) {
	ctx := context.Background()
	logger := slog.Default()

	mockClient := &mockHTTPClient{
		statusCode: http.StatusOK,
		body: `0:{"a":"$@1","f":"","b":"test"}
1:{"data":[{"office_name":"Viveknagar S.O","pincode":"560047"}],"success":true}`,
	}

	mockGeoManager := NewInstrumentedMockGeoJsonManager()
	mockGeoManager.SetMode(ModeNil) // Return nil for all GeoJSON queries

	provider := &BengaluruProvider{
		logger:     logger,
		geoManager: mockGeoManager,
		httpClient: mockClient,
	}

	entities, err := provider.GetEntities(ctx, 12.97, 77.59)
	if err != nil {
		t.Fatalf("GetEntities failed: %v", err)
	}

	// Find the Post Office entity
	var postOfficeEntity *models.Entity
	for i, e := range entities {
		if e.Name == "Post Office" {
			postOfficeEntity = &entities[i]
			break
		}
	}

	if postOfficeEntity == nil {
		t.Fatal("Post Office entity not found in response")
	}

	if !postOfficeEntity.IsAvailable {
		t.Error("Post Office entity should be available")
	}

	if len(postOfficeEntity.Attributes) == 0 {
		t.Error("Post Office entity should have attributes")
	}
}