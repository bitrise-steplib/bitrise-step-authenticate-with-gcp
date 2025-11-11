package step

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/v2/command"
)

func (s *Step) Run(config Config) (Result, error) {
	var keyPath string

	if config.ServiceAccountKey != "" {
		path, err := s.serviceAccountBasedAuthentication(config)
		if err != nil {
			return Result{}, fmt.Errorf("service account based authentication failed: %w", err)
		}
		keyPath = path
	} else if config.ClientConfig != "" {
		path, err := s.identityTokenBasedAuthentication(config)
		if err != nil {
			return Result{}, fmt.Errorf("identity token based authentication failed: %w", err)
		}
		keyPath = path
	} else {
		return Result{}, fmt.Errorf("no authentication method provided")
	}

	token, err := s.generateToken()
	if err != nil {
		return Result{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	s.logger.Printf("Access token generated")

	if config.DockerLogin {
		err = s.loginWithDocker(token, config.ArtifactRegistryLocations)
		if err != nil {
			return Result{}, fmt.Errorf("failed to login to Docker with GCP token: %w", err)
		}

		s.logger.Printf("Logged in with Docker")
	}

	return Result{
		AuthToken:       token,
		CredentialsPath: keyPath,
	}, nil
}

func save(data, pattern string) (string, error) {
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}

	if _, err := file.WriteString(data); err != nil {
		return "", fmt.Errorf("failed to write to temporary file: %w", err)
	}

	return file.Name(), nil
}

func (s *Step) generateToken() (string, error) {
	cmd := s.commandFactory.Create("gcloud", []string{"auth", "print-access-token", "--quiet"}, nil)
	output, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return "", err
	}

	return output, nil
}

func (s *Step) loginWithDocker(token string, locations []string) error {
	for _, location := range locations {
		cmd := s.commandFactory.Create("docker", []string{"login", "-u", "oauth2accesstoken", "--password-stdin", fmt.Sprintf("https://%s", location)}, &command.Opts{
			Stdin: strings.NewReader(token),
		})
		if output, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
			s.logger.Errorf("Docker login output: %s", output)
			return err
		}
	}

	return nil
}
