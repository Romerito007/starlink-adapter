package client

// Status is the normalized domain model for dish status.
type Status struct {
	DeviceID              string
	HardwareVersion       string
	SoftwareVersion       string
	UptimeSeconds         uint64
	UplinkThroughputBps   float32
	DownlinkThroughputBps float32
	PopPingDropRate       float32
	PopPingLatencyMs      float32
}

// Stats is the normalized domain model for dish history metrics.
type Stats struct {
	Current               uint64
	PopPingDropRate       []float32
	PopPingLatencyMs      []float32
	DownlinkThroughputBps []float32
	UplinkThroughputBps   []float32
}

// Location is the normalized domain model for dish coordinates.
type Location struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
	SigmaM    float64
	Source    string
}
