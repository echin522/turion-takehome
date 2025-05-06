package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Config variables pulled from user's environment. When service is deployed using
// k8s, these secrets would come from the service's configmap
type TelemetryGeneratorConfig struct {
	GroundStationEmulatorAddress string
}

func NewTelemetryGeneratorConfig() (*TelemetryGeneratorConfig, error) {
	gatewayName := strings.TrimSpace(os.Getenv("TELEMETRY_GATEWAY_SERVICE_NAME"))
	if gatewayName == "" {
		return nil, errors.New("env variable TELEMETRY_GATEWAY_SERVICE_NAME is empty")
	}

	groundStationEmulatorAddress := fmt.Sprintf("%s%s",
		gatewayName,
		strings.TrimSpace(os.Getenv("GROUND_STATION_EMULATOR_ADDRESS")),
	)
	if groundStationEmulatorAddress == "" {
		return nil, errors.New("env variable GROUND_STATION_EMULATOR_ADDRESS is empty")
	}

	return &TelemetryGeneratorConfig{
		GroundStationEmulatorAddress: (groundStationEmulatorAddress),
	}, nil
}
