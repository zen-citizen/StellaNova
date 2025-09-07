package config

import (
	"os"
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	tests := []struct {
		name         string
		envMap       map[string]string
		expectedPort string
		expectedGin  string
		expectedLog  string
		expectedFmt  string
	}{
		{
			name:         "all defaults when no env vars",
			envMap:       nil,
			expectedPort: "8080",
			expectedGin:  "release",
			expectedLog:  "info",
			expectedFmt:  "text",
		},
		{
			name: "custom port",
			envMap: map[string]string{
				"PORT": "3000",
			},
			expectedPort: "3000",
			expectedGin:  "release",
			expectedLog:  "info",
			expectedFmt:  "text",
		},
		{
			name: "custom log level",
			envMap: map[string]string{
				"LOG_LEVEL": "debug",
			},
			expectedPort: "8080",
			expectedGin:  "release",
			expectedLog:  "debug",
			expectedFmt:  "text",
		},
		{
			name: "multiple env vars",
			envMap: map[string]string{
				"LOG_LEVEL":  "error",
				"LOG_FORMAT": "json",
				"GIN_MODE":   "debug",
				"PORT":       "8080",
			},
			expectedPort: "8080",
			expectedGin:  "debug",
			expectedLog:  "error",
			expectedFmt:  "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env
			os.Clearenv()

			// Set environment variables from envMap
			for key, value := range tt.envMap {
				os.Setenv(key, value)
			}

			cfg := New()

			if cfg.Port != tt.expectedPort {
				t.Errorf("Port = %v, want %v", cfg.Port, tt.expectedPort)
			}
			if cfg.GinMode != tt.expectedGin {
				t.Errorf("GinMode = %v, want %v", cfg.GinMode, tt.expectedGin)
			}
			if cfg.LogLevel != tt.expectedLog {
				t.Errorf("LogLevel = %v, want %v", cfg.LogLevel, tt.expectedLog)
			}
			if cfg.LogFormat != tt.expectedFmt {
				t.Errorf("LogFormat = %v, want %v", cfg.LogFormat, tt.expectedFmt)
			}
		})
	}
}
