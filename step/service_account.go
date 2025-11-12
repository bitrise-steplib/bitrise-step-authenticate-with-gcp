package step

import (
	"encoding/json"
	"fmt"
)

func (s *Step) serviceAccountBasedAuthentication(config Config) (string, error) {
	s.logger.Println()
	s.logger.Infof("Performing service account based GCP authentication:")

	email, err := extractEmail(config.ServiceAccountKey)
	if err != nil {
		return "", fmt.Errorf("failed to extract email from service account key: %w", err)
	}

	keyPath, err := save(config.ServiceAccountKey, "service_account_key_*.json")
	if err != nil {
		return "", fmt.Errorf("failed to save service account key: %w", err)
	}

	if err := s.authenticate(email, keyPath); err != nil {
		return "", fmt.Errorf("failed to authenticate with service account: %w", err)
	}

	s.logger.Printf("GCP authentication successful")

	return keyPath, nil
}

func extractEmail(data string) (string, error) {
	var values map[string]string
	if err := json.Unmarshal([]byte(data), &values); err != nil {
		return "", err
	}

	email, ok := values["client_email"]
	if !ok {
		return "", fmt.Errorf("no client_email found")
	}

	return email, nil
}

func (s *Step) authenticate(email, keyPath string) error {
	cmd := s.commandFactory.Create("gcloud", []string{"auth", "activate-service-account", email, fmt.Sprintf("--key-file=%s", keyPath)}, nil)
	if output, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		s.logger.Errorf("GCP authentication output: %s", output)
		return err
	}

	return nil
}
