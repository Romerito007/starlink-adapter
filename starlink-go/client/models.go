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

// ClientDevice is the normalized domain model for a connected LAN/Wi-Fi client.
type ClientDevice struct {
	ClientID       uint32
	Name           string
	GivenName      string
	Domain         string
	MACAddress     string
	IPAddress      string
	IPv6Addresses  []string
	SignalStrength float32
	Interface      string
	InterfaceName  string
	Role           string
	DeviceID       string
}
