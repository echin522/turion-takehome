package turiondatapacket

// PayloadStats holds the min/max/avg for each field in TelemetryPayload.
// TODO: Define this in protobuf instead
type PayloadStats struct {
	MinTemperature float32 `json:"minTemperature"`
	MaxTemperature float32 `json:"maxTemperature"`
	AvgTemperature float32 `json:"avgTemperature"`

	MinBattery float32 `json:"minBattery"`
	MaxBattery float32 `json:"maxBattery"`
	AvgBattery float32 `json:"avgBattery"`

	MinAltitude float32 `json:"minAltitude"`
	MaxAltitude float32 `json:"maxAltitude"`
	AvgAltitude float32 `json:"avgAltitude"`

	MinSignal float32 `json:"minSignal"`
	MaxSignal float32 `json:"maxSignal"`
	AvgSignal float32 `json:"avgSignal"`
}
