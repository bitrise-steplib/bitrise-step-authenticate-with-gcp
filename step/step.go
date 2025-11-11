package step

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/log"
)

type Input struct {
	ServiceAccountKey         string          `env:"service_account_key"`
	ClientConfig              string          `env:"client_config"`
	Audience                  string          `env:"audience,required"`
	DockerLogin               bool            `env:"docker_login,opt[true,false]"`
	ArtifactRegistryLocations string          `env:"artifact_registry_locations"`
	BuildURL                  string          `env:"build_url,required"`
	BuildToken                stepconf.Secret `env:"build_api_token,required"`
	Verbose                   bool            `env:"verbose,opt[true,false]"`
}

type Config struct {
	ServiceAccountKey         string
	ClientConfig              string
	Audience                  string
	DockerLogin               bool
	ArtifactRegistryLocations []string
	BuildURL                  string
	BuildToken                stepconf.Secret
}

type Result struct {
	AuthToken       string
	CredentialsPath string
}

type Step struct {
	inputParser    stepconf.InputParser
	commandFactory command.Factory
	exporter       export.Exporter
	logger         log.Logger
}

func NewStep(
	inputParser stepconf.InputParser,
	commandFactory command.Factory,
	exporter export.Exporter,
	logger log.Logger,
) Step {
	return Step{
		inputParser:    inputParser,
		commandFactory: commandFactory,
		exporter:       exporter,
		logger:         logger,
	}
}

func (s *Step) ProcessConfig() (*Config, error) {
	var input Input
	err := s.inputParser.Parse(&input)
	if err != nil {
		return &Config{}, err
	}

	stepconf.Print(input)
	s.logger.EnableDebugLog(input.Verbose)

	if input.ServiceAccountKey != "" && input.ClientConfig != "" {
		return &Config{}, fmt.Errorf("only one authentication method can be used at a time (either Service Account or Identity Token)")
	}

	var locations []string
	for _, location := range strings.Split(input.ArtifactRegistryLocations, "\n") {
		if location == "" {
			continue
		}
		locations = append(locations, location)
	}

	if input.DockerLogin && len(locations) == 0 {
		return &Config{}, fmt.Errorf("no artifact registry locations specified for Docker login")
	}

	return &Config{
		ServiceAccountKey:         input.ServiceAccountKey,
		ClientConfig:              input.ClientConfig,
		Audience:                  input.Audience,
		DockerLogin:               input.DockerLogin,
		ArtifactRegistryLocations: locations,
		BuildURL:                  input.BuildURL,
		BuildToken:                input.BuildToken,
	}, nil
}
