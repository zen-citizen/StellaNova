package handlers

import (
	"backend/internal/models"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetEntitiesParameterValidation(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "missing lat parameter",
			queryParams:    "lng=77.5946&city=bengaluru",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "lat, lng, and city parameters are required",
		},
		{
			name:           "missing lng parameter",
			queryParams:    "lat=12.9716&city=bengaluru",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "lat, lng, and city parameters are required",
		},
		{
			name:           "invalid latitude - non-numeric",
			queryParams:    "lat=abc&lng=77.5946&city=bengaluru",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid latitude",
		},
		{
			name:           "invalid latitude - below minimum",
			queryParams:    "lat=-91&lng=77.5946&city=bengaluru",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid latitude",
		},
		{
			name:           "invalid latitude - above maximum",
			queryParams:    "lat=91&lng=77.5946&city=bengaluru",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid latitude",
		},
		{
			name:           "invalid longitude - non-numeric",
			queryParams:    "lat=12.9716&lng=xyz&city=bengaluru",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid longitude",
		},
		{
			name:           "invalid longitude - below minimum",
			queryParams:    "lat=12.9716&lng=-181&city=bengaluru",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid longitude",
		},
		{
			name:           "invalid longitude - above maximum",
			queryParams:    "lat=12.9716&lng=181&city=bengaluru",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid longitude",
		},
		{
			name:           "exactly at latitude minimum",
			queryParams:    "lat=-90&lng=77.5946&city=bengaluru",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "exactly at latitude maximum",
			queryParams:    "lat=90&lng=77.5946&city=bengaluru",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "exactly at longitude minimum",
			queryParams:    "lat=12.9716&lng=-180&city=bengaluru",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "exactly at longitude maximum",
			queryParams:    "lat=12.9716&lng=180&city=bengaluru",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			mockService := &mockEntitiesService{
				response: &models.EntitiesResponse{Entities: []models.Entity{}},
			}
			handler := NewEntitiesHandler(mockService, slog.Default())

			router := gin.New()
			router.GET("/entities", handler.GetEntities)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/entities?"+tt.queryParams, nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				if response["error"] != tt.expectedError {
					t.Errorf("expected error '%s', got '%s'", tt.expectedError, response["error"])
				}
			}
		})
	}
}

func TestGetEntitiesEmptyCityDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockEntitiesService{
		validateFunc: func(req *models.EntitiesRequest) {
			if req.City != "bengaluru" {
				t.Errorf("expected city to default to 'bengaluru', got '%s'", req.City)
			}
		},
		response: &models.EntitiesResponse{Entities: []models.Entity{}},
	}

	handler := NewEntitiesHandler(mockService, slog.Default())

	router := gin.New()
	router.GET("/entities", handler.GetEntities)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/entities?lat=12.9716&lng=77.5946", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGetEntitiesServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logBuffer strings.Builder
	logger := slog.New(slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	mockService := &mockEntitiesService{
		err: errors.New("something something service error"),
	}

	handler := NewEntitiesHandler(mockService, logger)

	router := gin.New()
	router.GET("/entities", handler.GetEntities)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/entities?lat=12.9716&lng=77.5946&city=bengaluru", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["error"] != "failed to get entities" {
		t.Errorf("expected error 'service error', got '%s'", response["error"])
	}

	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, "something something service error") {
		t.Error("expected log output to contain service error")
	}
}

type mockEntitiesService struct {
	response     *models.EntitiesResponse
	err          error
	validateFunc func(*models.EntitiesRequest)
}

func (m *mockEntitiesService) GetEntities(ctx context.Context, req *models.EntitiesRequest) (*models.EntitiesResponse, error) {
	if m.validateFunc != nil {
		m.validateFunc(req)
	}
	return m.response, m.err
}
