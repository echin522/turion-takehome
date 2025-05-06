// config/config_test.go
package config

import (
	"strings"
	"testing"
)

func TestNewTelemetryGeneratorConfig_Errors(t *testing.T) {
	tests := []struct {
		name            string
		telemetryEnv    string
		emulatorEnv     string
		wantErrContains string
	}{
		{
			name:            "both envs missing",
			telemetryEnv:    "",
			emulatorEnv:     "",
			wantErrContains: "TELEMETRY_GATEWAY_SERVICE_NAME is empty",
		},
		{
			name:            "telemetry only whitespace",
			telemetryEnv:    "   ",
			emulatorEnv:     "foo",
			wantErrContains: "TELEMETRY_GATEWAY_SERVICE_NAME is empty",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Clear and set envs
			t.Setenv("TELEMETRY_GATEWAY_SERVICE_NAME", tc.telemetryEnv)
			t.Setenv("GROUND_STATION_EMULATOR_ADDRESS", tc.emulatorEnv)

			cfg, err := NewTelemetryGeneratorConfig()
			if err == nil {
				t.Fatalf("expected error, got nil and cfg=%+v", cfg)
			}
			if !strings.Contains(err.Error(), tc.wantErrContains) {
				t.Errorf("error = %q; want to contain %q", err.Error(), tc.wantErrContains)
			}
		})
	}
}

func TestNewTelemetryGeneratorConfig_Success(t *testing.T) {
	tests := []struct {
		name         string
		telemetryEnv string
		emulatorEnv  string
		wantAddress  string
	}{
		{
			name:         "no slash in address",
			telemetryEnv: "gw",
			emulatorEnv:  ":8089",
			wantAddress:  "gw:8089",
		},
		{
			name:         "with whitespace to trim",
			telemetryEnv: "  my-gateway  ",
			emulatorEnv:  "  /api/v1/data  ",
			wantAddress:  "my-gateway/api/v1/data",
		},
		{
			name:         "empty emulator part yields just gateway",
			telemetryEnv: "server",
			emulatorEnv:  "   ",
			wantAddress:  "server",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TELEMETRY_GATEWAY_SERVICE_NAME", tc.telemetryEnv)
			t.Setenv("GROUND_STATION_EMULATOR_ADDRESS", tc.emulatorEnv)

			cfg, err := NewTelemetryGeneratorConfig()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg.GroundStationEmulatorAddress != tc.wantAddress {
				t.Errorf("GroundStationEmulatorAddress = %q; want %q",
					cfg.GroundStationEmulatorAddress, tc.wantAddress)
			}
		})
	}
}
