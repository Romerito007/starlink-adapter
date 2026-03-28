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

// StatusDetailed is the normalized detailed status snapshot for operational use.
type StatusDetailed struct {
	DeviceID                                    string
	HardwareVersion                             string
	SoftwareVersion                             string
	UptimeSeconds                               uint64
	Ipv4WanAddress                              string
	Ipv6WanAddresses                            []string
	PingLatencyMs                               float32
	PingDropRate                                float32
	PingDropRate5m                              float32
	DishPingLatencyMs                           float32
	DishPingDropRate                            float32
	DishPingDropRate5m                          float32
	PopPingLatencyMs                            float32
	PopPingDropRate                             float32
	PopPingDropRate5m                           float32
	PopIpv6PingLatencyMs                        float32
	PopIpv6PingDropRate                         float32
	PopIpv6PingDropRate5m                       float32
	SecsSinceLastPublicIpv4Change               uint32
	DishID                                      string
	UtcNs                                       int64
	DishDisablementCode                         string
	CalibrationPartitionsState                  string
	SetupRequirementState                       string
	SoftwareUpdateState                         string
	SoftwareUpdateRunningVersion                string
	SoftwareUpdateSecondsSinceGetTargetVersions float32
	PoeState                                    string
	PoePower                                    float32
	PoeVin                                      float32
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

// WifiConfigSnapshot is the normalized public view of wifi_get_config.
type WifiConfigSnapshot struct {
	CountryCode     string
	SetupComplete   bool
	MacWan          string
	MacLan          string
	BootCount       int32
	Incarnation     uint64
	WanHostDscpMark int32
	Networks        []WifiNetwork
}

// WifiNetwork is the normalized public network configuration model.
type WifiNetwork struct {
	Ipv4                       string
	Domain                     string
	Dhcpv4Start                uint32
	Dhcpv4End                  uint32
	Dhcpv4LeaseDurationSeconds uint32
	Vlan                       uint32
	BasicServiceSets           []WifiBasicServiceSet
}

// WifiBasicServiceSet is the normalized public BSS model.
type WifiBasicServiceSet struct {
	Bssid         string
	Ssid          string
	Band          string
	InterfaceName string
}

// NetworkInterfaceSnapshot is the normalized public view of a network interface.
type NetworkInterfaceSnapshot struct {
	Name          string
	Up            bool
	MacAddress    string
	Ipv4Addresses []string
	Ipv6Addresses []string
	RxStats       InterfaceTrafficStats
	TxStats       InterfaceTrafficStats
	Ethernet      *InterfaceEthernetInfo
	Wifi          *InterfaceWifiInfo
	Bridge        *InterfaceBridgeInfo
}

// InterfaceTrafficStats is the normalized traffic counters for an interface.
type InterfaceTrafficStats struct {
	Bytes   uint64
	Packets uint64
}

// InterfaceEthernetInfo is the normalized ethernet-specific view.
type InterfaceEthernetInfo struct {
	LinkDetected      bool
	SpeedMbps         uint32
	AutonegotiationOn bool
	Duplex            string
}

// InterfaceWifiInfo is the normalized wifi-specific view.
type InterfaceWifiInfo struct {
	Channel     uint32
	LinkQuality float64
}

// InterfaceBridgeInfo is the normalized bridge-specific view.
type InterfaceBridgeInfo struct {
	MemberNames []string
}

// RadioStat is the normalized public view of per-band radio health.
type RadioStat struct {
	Band          string
	RxStats       RadioTrafficStats
	TxStats       RadioTrafficStats
	ThermalStatus RadioThermalStatus
	AntennaStatus RadioAntennaStatus
}

// RadioTrafficStats is the normalized packet/error counters for radio traffic.
type RadioTrafficStats struct {
	Packets     uint64
	FrameErrors uint64
}

// RadioThermalStatus is the normalized thermal view for a radio.
type RadioThermalStatus struct {
	Temp      float64
	DutyCycle uint32
}

// RadioAntennaStatus is the normalized RSSI view for up to four antenna chains.
type RadioAntennaStatus struct {
	Rssi1 float32
	Rssi2 float32
	Rssi3 float32
	Rssi4 float32
}
