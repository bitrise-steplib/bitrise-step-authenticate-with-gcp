package step

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-steplib/bitrise-step-get-identity-token/api"
)

func (s *Step) identityTokenBasedAuthentication(config Config) (string, error) {
	identityToken, err := s.identityToken(config.BuildURL, config.BuildToken, config.Audience)
	if err != nil {
		return "", err
	}

	s.logger.Println()
	s.logger.Printf("Identity token fetched.")

	s.logger.Println()
	s.logger.Infof("Performing identity token based GCP authentication:")

	identityTokenPath, err := save(identityToken, "identity_token*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to save identity token: %w", err)
	}

	updatedClientConfig, err := updateClientConfigWithToken(config.ClientConfig, identityTokenPath)
	if err != nil {
		return "", fmt.Errorf("failed to update client config: %w", err)
	}

	configPath, err := save(updatedClientConfig, "client_config_*.json")
	if err != nil {
		return "", fmt.Errorf("failed to save client config: %w", err)
	}

	if err := s.login(configPath); err != nil {
		return "", fmt.Errorf("failed to authenticate with identity token: %w", err)
	}

	s.logger.Printf("GCP authentication successful")

	return configPath, nil
}

func (s *Step) identityToken(url string, token stepconf.Secret, aud string) (string, error) {
	client := api.NewDefaultAPIClient(url, token, s.logger)

	parameter := api.GetIdentityTokenParameter{
		Audience: aud,
	}
	response, err := client.GetIdentityToken(parameter)
	if err != nil {
		return "", err
	}

	return response.Token, nil
}

func updateClientConfigWithToken(clientConfig, identityTokenPath string) (string, error) {
	var values map[string]interface{}
	if err := json.Unmarshal([]byte(clientConfig), &values); err != nil {
		return "", fmt.Errorf("failed to unmarshal client config: %w", err)
	}
	values["credential_source"] = map[string]interface{}{
		"file": identityTokenPath,
		"format": map[string]interface{}{
			"type": "text",
		},
	}

	updatedConfigBytes, err := json.Marshal(values)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated client config: %w", err)
	}

	return string(updatedConfigBytes), nil
}

func (s *Step) login(configPath string) error {
	cmd := s.commandFactory.Create("gcloud", []string{"auth", "login", fmt.Sprintf("--cred-file=%s", configPath)}, nil)
	if output, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		s.logger.Errorf("GCP authentication output: %s", output)
		return err
	}

	return nil
}
