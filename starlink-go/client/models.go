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
	MacAddress                   string
	IpAddress                    string
	Interface                    string
	InterfaceName                string
	UpstreamMacAddress           string
	Role                         string
	SignalStrength               float32
	Snr                          float32
	ChannelWidth                 uint32
	Mode                         string
	Blocked                      bool
	DhcpLeaseActive              bool
	DhcpLeaseRenewed             bool
	DhcpLeaseFound               bool
	SecondsUntilDhcpLeaseExpires uint32
	AssociatedTimeSeconds        uint32
	NoDataIdleSeconds            uint32
	HopsFromController           uint32
	ClientID                     uint32
	CaptiveClientID              uint32
	UploadMb                     float32
	DownloadMb                   float32
	RxStatsValid                 bool
	TxStatsValid                 bool
	RxBytes                      uint64
	TxBytes                      uint64
	RxNss                        int32
	TxNss                        int32
	RxMcs                        uint32
	TxMcs                        uint32
	RxBandwidth                  uint32
	TxBandwidth                  uint32
	RxGuardNs                    uint32
	TxGuardNs                    uint32
	RxRateMbps                   uint32
	TxRateMbps                   uint32
	RxPhyMode                    uint32
	TxPhyMode                    uint32
	RxRateMbpsLast15s            float32
	TxRateMbpsLast15s            float32
	RxRateMbpsLast1mAvg          float32
	TxRateMbpsLast30s            float32
	Name                         string
	GivenName                    string
	Domain                       string
	Ipv6Addresses                []string
}

// DhcpLease is the normalized domain model for a DHCP lease entry.
type DhcpLease struct {
	IpAddress   string
	MacAddress  string
	Hostname    string
	ExpiresTime string
	Active      bool
	ClientID    uint32
	Domain      string
}
