package turiondatapacket

// CCSDSPrimaryHeader represents the 6‐byte primary header
type CCSDSPrimaryHeader struct {
	PacketID      uint16 `json:"packetId"`      // Version/type/APID
	PacketSeqCtrl uint16 `json:"packetSeqCtrl"` // SeqFlags/SeqCount
	PacketLength  uint16 `json:"packetLength"`  // Packet length minus 7
}

// CCSDSSecondaryHeader represents the 10‐byte secondary header
type CCSDSSecondaryHeader struct {
	Timestamp   uint64 `json:"timestamp"`   // Unix timestamp (seconds)
	SubsystemID uint16 `json:"subsystemId"` // e.g. power, thermal subsystem
}

// TelemetryPayload is the actual sensor data
type TelemetryPayload struct {
	Temperature float32 `json:"temperature"` // °C
	Battery     float32 `json:"battery"`     // %
	Altitude    float32 `json:"altitude"`    // km
	Signal      float32 `json:"signal"`      // dB
}

// TurionDataPacket is the full packet: headers + payload
type TurionDataPacket struct {
	CCSDSPrimaryHeader   CCSDSPrimaryHeader   `json:"ccsdsPrimaryHeader"`
	CCSDSSecondaryHeader CCSDSSecondaryHeader `json:"ccsdsSecondaryHeader"`
	TelemetryPayload     TelemetryPayload     `json:"telemetryPayload"`
}

const (
	APID           = 0x01
	PACKET_VERSION = 0x0
	PACKET_TYPE    = 0x0    // 0 = TM (telemetry)
	SEC_HDR_FLAG   = 0x1    // Secondary header present
	SEQ_FLAGS      = 0x3    // Standalone packet
	SUBSYSTEM_ID   = 0x0001 // Main bus telemetry
)

// DetectAnomalies returns one Anomaly per field whose value is outside its normal range.
func (dp TurionDataPacket) DetectAnomalies() []Anomaly {
	ts := dp.CCSDSSecondaryHeader.Timestamp
	var anomalies []Anomaly

	v := dp.TelemetryPayload

	// Temperature: normal 20–30; anomaly if >35
	if v.Temperature > 35.0 {
		anomalies = append(
			anomalies,
			Anomaly{Value: v.Temperature, Timestamp: ts, Field: Field_TEMPERATURE},
		)
	}

	// Battery: normal 70–100; anomaly if <40
	if v.Battery < 40.0 {
		anomalies = append(
			anomalies,
			Anomaly{Value: v.Battery, Timestamp: ts, Field: Field_BATTERY},
		)
	}

	// Altitude: normal 500–550; anomaly if <400
	if v.Altitude < 400.0 {
		anomalies = append(
			anomalies,
			Anomaly{Value: v.Altitude, Timestamp: ts, Field: Field_ALTITUDE},
		)
	}

	// Signal: normal -60…-40; anomaly if < -80
	if v.Signal < -80.0 {
		anomalies = append(
			anomalies,
			Anomaly{Value: v.Signal, Timestamp: ts, Field: Field_SIGNAL},
		)
	}

	return anomalies
}
