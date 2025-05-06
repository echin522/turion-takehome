package config

import (
	"errors"
	"os"
	"strings"
)

// TelemetryAPIConfig is the env config for the telemetry API service
type TelemetryAPIConfig struct {
	PGHostURL string
}

func NewTelemetryAPIConfig() (*TelemetryAPIConfig, error) {
	pgHostURL := strings.TrimSpace(os.Getenv("PG_HOST_URL"))
	if pgHostURL == "" {
		return nil, errors.New("env variable PG_HOST_URL is empty")
	}

	return &TelemetryAPIConfig{
		PGHostURL: pgHostURL,
	}, nil
}
