package config

import (
	"errors"
	"os"
	"strings"
)

// Config variables pulled from user's environment. When service is deployed using
// k8s, these secrets would come from the service's configmap
type TelemetryGatewayConfig struct {
	PGHostURL                    string
	GroundStationEmulatorAddress string
	TelemetryAPIServerURL        string
}

func NewTelemetryGatewayConfig() (*TelemetryGatewayConfig, error) {
	pgHostURL := strings.TrimSpace(os.Getenv("PG_HOST_URL"))
	if pgHostURL == "" {
		return nil, errors.New("env variable PG_HOST_URL is empty")
	}

	telemetryAPIServerURL := strings.TrimSpace(os.Getenv("TELEMETRY_API_SERVER_URL"))
	if pgHostURL == "" {
		return nil, errors.New("env variable TELEMETRY_API_SERVER_URL is empty")
	}

	groundStationEmulatorAddress := strings.TrimSpace(os.Getenv("GROUND_STATION_EMULATOR_ADDRESS"))
	if groundStationEmulatorAddress == "" {
		return nil, errors.New("env variable GROUND_STATION_EMULATOR_ADDRESS is empty")
	}

	return &TelemetryGatewayConfig{
		PGHostURL:                    pgHostURL,
		TelemetryAPIServerURL:        telemetryAPIServerURL,
		GroundStationEmulatorAddress: groundStationEmulatorAddress,
	}, nil
}
