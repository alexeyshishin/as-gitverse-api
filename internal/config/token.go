package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	tokenEnvVar     = "GITVERSE_TOKEN"
	credentialsFile = "credentials"
	credentialsDir  = "gitverse"
)

func ResolveToken(flagValue string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}

	if envValue := os.Getenv(tokenEnvVar); envValue != "" {
		return envValue, nil
	}

	path, err := credentialsPath()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("no token: pass --token, set %s, or write one to %s", tokenEnvVar, path)
		}
		return "", fmt.Errorf("read credentials file: %w", err)
	}

	token := strings.TrimSpace(string(data))
	if token == "" {
		return "", fmt.Errorf("credentials file %s is empty", path)
	}
	return token, nil
}

func WriteToken(token string) error {
	path, err := credentialsPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	return os.WriteFile(path, []byte(token+"\n"), 0o600)
}

func credentialsPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}
	return filepath.Join(configDir, credentialsDir, credentialsFile), nil
}
