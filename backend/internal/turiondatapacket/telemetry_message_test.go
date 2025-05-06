package turiondatapacket

import (
	"reflect"
	"testing"
	"time"
)

// helper to make a packet quickly
func makePacket(temp, batt, alt, sig float32, ts uint64) TurionDataPacket {
	return TurionDataPacket{
		CCSDSSecondaryHeader: CCSDSSecondaryHeader{Timestamp: ts},
		TelemetryPayload: TelemetryPayload{
			Temperature: temp,
			Battery:     batt,
			Altitude:    alt,
			Signal:      sig,
		},
	}
}

func TestDetectAnomalies(t *testing.T) {
	now := uint64(time.Now().Unix())

	tests := []struct {
		name string
		pkt  TurionDataPacket
		want []Anomaly
	}{
		{
			name: "no anomalies at perfect normals",
			pkt:  makePacket(25, 85, 525, -50, now),
			want: nil,
		},
		{
			name: "temperature anomaly only",
			pkt:  makePacket(36, 85, 525, -50, now),
			want: []Anomaly{{Value: 36, Timestamp: now, Field: Field_TEMPERATURE}},
		},
		{
			name: "battery anomaly only",
			pkt:  makePacket(25, 39, 525, -50, now),
			want: []Anomaly{{Value: 39, Timestamp: now, Field: Field_BATTERY}},
		},
		{
			name: "altitude anomaly only",
			pkt:  makePacket(25, 85, 399, -50, now),
			want: []Anomaly{{Value: 399, Timestamp: now, Field: Field_ALTITUDE}},
		},
		{
			name: "signal anomaly only",
			pkt:  makePacket(25, 85, 525, -81, now),
			want: []Anomaly{{Value: -81, Timestamp: now, Field: Field_SIGNAL}},
		},
		{
			name: "multiple anomalies",
			pkt:  makePacket(50, 20, 350, -90, now),
			want: []Anomaly{
				{Value: 50, Timestamp: now, Field: Field_TEMPERATURE},
				{Value: 20, Timestamp: now, Field: Field_BATTERY},
				{Value: 350, Timestamp: now, Field: Field_ALTITUDE},
				{Value: -90, Timestamp: now, Field: Field_SIGNAL},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pkt.DetectAnomalies()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectAnomalies() = %v; want %v", got, tt.want)
			}
		})
	}
}
