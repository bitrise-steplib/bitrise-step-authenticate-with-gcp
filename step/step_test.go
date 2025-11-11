package step

import (
	"testing"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-authenticate-wth-gcp/mocks"
	"github.com/stretchr/testify/assert"
)

func TestConfigParsing(t *testing.T) {
	tests := []struct {
		name           string
		env            map[string]string
		expectedConfig Config
		expectError    bool
	}{
		{
			name: "valid config with service account key",
			env: map[string]string{
				"service_account_key":         "service-account-key",
				"client_config":               "",
				"audience":                    "audience",
				"docker_login":                "true",
				"artifact_registry_locations": "location-1\nlocation-2",
				"build_url":                   "build-url",
				"build_api_token":             "build-token",
				"verbose":                     "true",
			},
			expectedConfig: Config{
				ServiceAccountKey:         "service-account-key",
				ClientConfig:              "",
				Audience:                  "audience",
				DockerLogin:               true,
				ArtifactRegistryLocations: []string{"location-1", "location-2"},
				BuildURL:                  "build-url",
				BuildToken:                "build-token",
			},
			expectError: false,
		},
		{
			name: "valid config with client config",
			env: map[string]string{
				"service_account_key":         "",
				"client_config":               "client-config",
				"audience":                    "audience",
				"docker_login":                "false",
				"artifact_registry_locations": "",
				"build_url":                   "build-url",
				"build_api_token":             "build-token",
				"verbose":                     "false",
			},
			expectedConfig: Config{
				ServiceAccountKey:         "",
				ClientConfig:              "client-config",
				Audience:                  "audience",
				DockerLogin:               false,
				ArtifactRegistryLocations: nil,
				BuildURL:                  "build-url",
				BuildToken:                "build-token",
			},
			expectError: false,
		},
		{
			name: "error when both service account key and client config are set",
			env: map[string]string{
				"service_account_key":         "service-account-key",
				"client_config":               "client-config",
				"audience":                    "audience",
				"docker_login":                "false",
				"artifact_registry_locations": "",
				"build_url":                   "build-url",
				"build_api_token":             "build-token",
				"verbose":                     "false",
			},
			expectedConfig: Config{},
			expectError:    true,
		},
		{
			name: "error when docker login is true but no locations",
			env: map[string]string{
				"service_account_key":         "service-account-key",
				"client_config":               "",
				"audience":                    "audience",
				"docker_login":                "true",
				"artifact_registry_locations": "",
				"build_url":                   "build-url",
				"build_api_token":             "build-token",
				"verbose":                     "false",
			},
			expectedConfig: Config{},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEnvRepository := mocks.NewRepository(t)
			for k, v := range tt.env {
				mockEnvRepository.On("Get", k).Return(v)
			}

			inputParser := stepconf.NewInputParser(mockEnvRepository)
			mockFactory := mocks.NewFactory(t)
			exporter := export.NewExporter(mocks.NewFactory(t))
			sut := NewStep(inputParser, mockFactory, exporter, log.NewLogger())

			receivedConfig, err := sut.ProcessConfig()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedConfig, *receivedConfig)
			}

			mockEnvRepository.AssertExpectations(t)
		})
	}
}
