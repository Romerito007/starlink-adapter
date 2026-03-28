package client

import (
	"math"
	"sort"
	"strings"

	pb "github.com/Romerito007/starlink-adapter/starlink-go/proto/gen/spacex/api/device"
)

func mapStatus(in *pb.DishGetStatusResponse) *Status {
	if in == nil {
		return &Status{}
	}

	deviceInfo := in.GetDeviceInfo()
	deviceState := in.GetDeviceState()

	return &Status{
		DeviceID:              deviceInfo.GetId(),
		HardwareVersion:       deviceInfo.GetHardwareVersion(),
		SoftwareVersion:       deviceInfo.GetSoftwareVersion(),
		UptimeSeconds:         deviceState.GetUptimeS(),
		UplinkThroughputBps:   in.GetUplinkThroughputBps(),
		DownlinkThroughputBps: in.GetDownlinkThroughputBps(),
		PopPingDropRate:       in.GetPopPingDropRate(),
		PopPingLatencyMs:      in.GetPopPingLatencyMs(),
	}
}

func mapStatusFromWifi(in *pb.WifiGetStatusResponse) *Status {
	if in == nil {
		return &Status{}
	}

	deviceInfo := in.GetDeviceInfo()
	deviceState := in.GetDeviceState()

	return &Status{
		DeviceID:              deviceInfo.GetId(),
		HardwareVersion:       deviceInfo.GetHardwareVersion(),
		SoftwareVersion:       deviceInfo.GetSoftwareVersion(),
		UptimeSeconds:         deviceState.GetUptimeS(),
		UplinkThroughputBps:   0,
		DownlinkThroughputBps: 0,
		PopPingDropRate:       in.GetPopPingDropRate(),
		PopPingLatencyMs:      in.GetPopPingLatencyMs(),
	}
}

func mapStats(in *pb.DishGetHistoryResponse) *Stats {
	if in == nil {
		return &Stats{}
	}

	return &Stats{
		Current:               in.GetCurrent(),
		PopPingDropRate:       append([]float32(nil), in.GetPopPingDropRate()...),
		PopPingLatencyMs:      append([]float32(nil), in.GetPopPingLatencyMs()...),
		DownlinkThroughputBps: append([]float32(nil), in.GetDownlinkThroughputBps()...),
		UplinkThroughputBps:   append([]float32(nil), in.GetUplinkThroughputBps()...),
	}
}

func mapLocation(in *pb.GetLocationResponse) *Location {
	if in == nil {
		return &Location{}
	}

	lla := in.GetLla()
	if lla == nil {
		return &Location{
			SigmaM: in.GetSigmaM(),
			Source: in.GetSource().String(),
		}
	}

	return &Location{
		Latitude:  lla.GetLat(),
		Longitude: lla.GetLon(),
		Altitude:  lla.GetAlt(),
		SigmaM:    in.GetSigmaM(),
		Source:    in.GetSource().String(),
	}
}

func mapConnectedClients(in []*pb.WifiClient) []ClientDevice {
	if len(in) == 0 {
		return []ClientDevice{}
	}

	out := make([]ClientDevice, 0, len(in))
	for _, c := range in {
		if c == nil {
			continue
		}

		ipv6 := append([]string{}, c.GetIpv6Addresses()...)
		sort.Strings(ipv6)
		rxStats := c.GetRxStats()
		txStats := c.GetTxStats()
		rates := mapRecentRates(rxStats, txStats)

		out = append(out, ClientDevice{
			MacAddress:                   c.GetMacAddress(),
			IpAddress:                    c.GetIpAddress(),
			Interface:                    normalizeEnum(c.GetIface().String(), ""),
			InterfaceName:                c.GetIfaceName(),
			UpstreamMacAddress:           c.GetUpstreamMacAddress(),
			Role:                         normalizeEnum(c.GetRole().String(), "ROLE_"),
			SignalStrength:               c.GetSignalStrength(),
			Snr:                          c.GetSnr(),
			ChannelWidth:                 c.GetChannelWidth(),
			Mode:                         c.GetModeStr(),
			Blocked:                      c.GetBlocked(),
			DhcpLeaseActive:              c.GetDhcpLeaseActive(),
			DhcpLeaseRenewed:             c.GetDhcpLeaseRenewed(),
			DhcpLeaseFound:               false,
			SecondsUntilDhcpLeaseExpires: 0,
			AssociatedTimeSeconds:        c.GetAssociatedTimeS(),
			NoDataIdleSeconds:            c.GetNoDataIdleS(),
			HopsFromController:           c.GetHopsFromController(),
			ClientID:                     c.GetClientId(),
			CaptiveClientID:              0,
			UploadMb:                     0,
			DownloadMb:                   0,
			RxStatsValid:                 rxStats != nil,
			TxStatsValid:                 txStats != nil,
			RxBytes:                      rxStats.GetBytes(),
			TxBytes:                      txStats.GetBytes(),
			RxNss:                        rxStats.GetNss(),
			TxNss:                        txStats.GetNss(),
			RxMcs:                        rxStats.GetMcs(),
			TxMcs:                        txStats.GetMcs(),
			RxBandwidth:                  rxStats.GetBandwidth(),
			TxBandwidth:                  txStats.GetBandwidth(),
			RxGuardNs:                    rxStats.GetGuardNs(),
			TxGuardNs:                    txStats.GetGuardNs(),
			RxRateMbps:                   rates.rxRateMbps,
			TxRateMbps:                   rates.txRateMbps,
			RxPhyMode:                    rxStats.GetPhyMode(),
			TxPhyMode:                    txStats.GetPhyMode(),
			RxRateMbpsLast15s:            rates.rxRateMbpsLast15s,
			TxRateMbpsLast15s:            rates.txRateMbpsLast15s,
			RxRateMbpsLast1mAvg:          rates.rxRateMbpsLast1mAvg,
			TxRateMbpsLast30s:            rates.txRateMbpsLast30s,
			Name:                         c.GetName(),
			GivenName:                    c.GetGivenName(),
			Domain:                       c.GetDomain(),
			Ipv6Addresses:                ipv6,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		leftIface := strings.ToLower(out[i].Interface)
		rightIface := strings.ToLower(out[j].Interface)
		if leftIface != rightIface {
			return leftIface < rightIface
		}

		leftMAC := strings.ToLower(out[i].MacAddress)
		rightMAC := strings.ToLower(out[j].MacAddress)
		if leftMAC != rightMAC {
			return leftMAC < rightMAC
		}

		leftName := strings.ToLower(out[i].Name)
		rightName := strings.ToLower(out[j].Name)
		return leftName < rightName
	})

	return out
}

func mapDhcpLeases(in []*pb.DhcpServer) []DhcpLease {
	if len(in) == 0 {
		return []DhcpLease{}
	}

	out := make([]DhcpLease, 0)
	for _, server := range in {
		if server == nil {
			continue
		}

		domain := server.GetDomain()
		for _, lease := range server.GetLeases() {
			if lease == nil {
				continue
			}

			out = append(out, DhcpLease{
				IpAddress:   lease.GetIpAddress(),
				MacAddress:  lease.GetMacAddress(),
				Hostname:    lease.GetHostname(),
				ExpiresTime: lease.GetExpiresTime(),
				Active:      lease.GetActive(),
				ClientID:    lease.GetClientId(),
				Domain:      domain,
			})
		}
	}

	sort.Slice(out, func(i, j int) bool {
		leftDomain := strings.ToLower(out[i].Domain)
		rightDomain := strings.ToLower(out[j].Domain)
		if leftDomain != rightDomain {
			return leftDomain < rightDomain
		}

		leftIP := strings.ToLower(out[i].IpAddress)
		rightIP := strings.ToLower(out[j].IpAddress)
		if leftIP != rightIP {
			return leftIP < rightIP
		}

		leftMAC := strings.ToLower(out[i].MacAddress)
		rightMAC := strings.ToLower(out[j].MacAddress)
		return leftMAC < rightMAC
	})

	return out
}

func mapWifiConfigSnapshot(in *pb.WifiGetConfigResponse) *WifiConfigSnapshot {
	if in == nil || in.GetWifiConfig() == nil {
		return &WifiConfigSnapshot{
			Networks: []WifiNetwork{},
		}
	}

	cfg := in.GetWifiConfig()
	networks := mapWifiNetworks(cfg.GetNetworks())

	return &WifiConfigSnapshot{
		CountryCode:     cfg.GetCountryCode(),
		SetupComplete:   cfg.GetSetupComplete(),
		MacWan:          cfg.GetMacWan(),
		MacLan:          cfg.GetMacLan(),
		BootCount:       cfg.GetBootCount(),
		Incarnation:     cfg.GetIncarnation(),
		WanHostDscpMark: cfg.GetWanHostDscpMark(),
		Networks:        networks,
	}
}

func mapWifiNetworks(in []*pb.WifiConfig_Network) []WifiNetwork {
	if len(in) == 0 {
		return []WifiNetwork{}
	}

	out := make([]WifiNetwork, 0, len(in))
	for _, n := range in {
		if n == nil {
			continue
		}

		out = append(out, WifiNetwork{
			Ipv4:                       n.GetIpv4(),
			Domain:                     n.GetDomain(),
			Dhcpv4Start:                n.GetDhcpv4Start(),
			Dhcpv4End:                  0,
			Dhcpv4LeaseDurationSeconds: n.GetDhcpv4LeaseDurationS(),
			Vlan:                       n.GetVlan(),
			BasicServiceSets:           mapWifiBasicServiceSets(n.GetBasicServiceSets()),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		leftDomain := strings.ToLower(out[i].Domain)
		rightDomain := strings.ToLower(out[j].Domain)
		if leftDomain != rightDomain {
			return leftDomain < rightDomain
		}

		leftIP := strings.ToLower(out[i].Ipv4)
		rightIP := strings.ToLower(out[j].Ipv4)
		if leftIP != rightIP {
			return leftIP < rightIP
		}

		if out[i].Vlan != out[j].Vlan {
			return out[i].Vlan < out[j].Vlan
		}
		return len(out[i].BasicServiceSets) < len(out[j].BasicServiceSets)
	})

	return out
}

func mapWifiBasicServiceSets(in []*pb.WifiConfig_BasicServiceSet) []WifiBasicServiceSet {
	if len(in) == 0 {
		return []WifiBasicServiceSet{}
	}

	out := make([]WifiBasicServiceSet, 0, len(in))
	for _, bss := range in {
		if bss == nil {
			continue
		}

		out = append(out, WifiBasicServiceSet{
			Bssid:         bss.GetBssid(),
			Ssid:          bss.GetSsid(),
			Band:          normalizeEnum(bss.GetBand().String(), ""),
			InterfaceName: bss.GetIfaceName(),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		leftSSID := strings.ToLower(out[i].Ssid)
		rightSSID := strings.ToLower(out[j].Ssid)
		if leftSSID != rightSSID {
			return leftSSID < rightSSID
		}

		leftBSSID := strings.ToLower(out[i].Bssid)
		rightBSSID := strings.ToLower(out[j].Bssid)
		if leftBSSID != rightBSSID {
			return leftBSSID < rightBSSID
		}

		leftIface := strings.ToLower(out[i].InterfaceName)
		rightIface := strings.ToLower(out[j].InterfaceName)
		return leftIface < rightIface
	})

	return out
}

func mapNetworkInterfaces(in []*pb.NetworkInterface) []NetworkInterfaceSnapshot {
	if len(in) == 0 {
		return []NetworkInterfaceSnapshot{}
	}

	out := make([]NetworkInterfaceSnapshot, 0, len(in))
	for _, iface := range in {
		if iface == nil {
			continue
		}

		ipv4 := append([]string{}, iface.GetIpv4Addresses()...)
		sort.Strings(ipv4)
		ipv6 := append([]string{}, iface.GetIpv6Addresses()...)
		sort.Strings(ipv6)

		out = append(out, NetworkInterfaceSnapshot{
			Name:          iface.GetName(),
			Up:            iface.GetUp(),
			MacAddress:    iface.GetMacAddress(),
			Ipv4Addresses: ipv4,
			Ipv6Addresses: ipv6,
			RxStats: InterfaceTrafficStats{
				Bytes:   iface.GetRxStats().GetBytes(),
				Packets: iface.GetRxStats().GetPackets(),
			},
			TxStats: InterfaceTrafficStats{
				Bytes:   iface.GetTxStats().GetBytes(),
				Packets: iface.GetTxStats().GetPackets(),
			},
			Ethernet: mapInterfaceEthernet(iface.GetEthernet()),
			Wifi:     mapInterfaceWifi(iface.GetWifi()),
			Bridge:   mapInterfaceBridge(iface.GetBridge()),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})

	return out
}

func mapInterfaceEthernet(in *pb.EthernetNetworkInterface) *InterfaceEthernetInfo {
	if in == nil {
		return nil
	}

	return &InterfaceEthernetInfo{
		LinkDetected:      in.GetLinkDetected(),
		SpeedMbps:         in.GetSpeedMbps(),
		AutonegotiationOn: in.GetAutonegotiationOn(),
		Duplex:            normalizeEnum(in.GetDuplex().String(), ""),
	}
}

func mapInterfaceWifi(in *pb.WifiNetworkInterface) *InterfaceWifiInfo {
	if in == nil {
		return nil
	}

	return &InterfaceWifiInfo{
		Channel:     in.GetChannel(),
		LinkQuality: in.GetLinkQuality(),
	}
}

func mapInterfaceBridge(in *pb.BridgeNetworkInterface) *InterfaceBridgeInfo {
	if in == nil {
		return nil
	}

	members := append([]string{}, in.GetMemberNames()...)
	sort.Strings(members)

	return &InterfaceBridgeInfo{
		MemberNames: members,
	}
}

func mapRadioStats(in []*pb.RadioStats) []RadioStat {
	if len(in) == 0 {
		return []RadioStat{}
	}

	out := make([]RadioStat, 0, len(in))
	for _, radio := range in {
		if radio == nil {
			continue
		}

		rxStats := radio.GetRxStats()
		txStats := radio.GetTxStats()
		thermal := radio.GetThermalStatus()
		antenna := radio.GetAntennaStatus()

		out = append(out, RadioStat{
			Band: normalizeEnum(radio.GetBand().String(), ""),
			RxStats: RadioTrafficStats{
				Packets:     rxStats.GetPackets(),
				FrameErrors: rxStats.GetFrameErrors(),
			},
			TxStats: RadioTrafficStats{
				Packets:     txStats.GetPackets(),
				FrameErrors: 0,
			},
			ThermalStatus: RadioThermalStatus{
				Temp:      sanitizeFloat64(thermal.GetTemp2()),
				DutyCycle: thermal.GetDutyCycle(),
			},
			AntennaStatus: RadioAntennaStatus{
				Rssi1: sanitizeFloat32(antenna.GetRssi1()),
				Rssi2: sanitizeFloat32(antenna.GetRssi2()),
				Rssi3: sanitizeFloat32(antenna.GetRssi3()),
				Rssi4: sanitizeFloat32(antenna.GetRssi4()),
			},
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Band) < strings.ToLower(out[j].Band)
	})

	return out
}

func sanitizeFloat32(v float32) float32 {
	if math.IsNaN(float64(v)) {
		return 0
	}
	return v
}

func sanitizeFloat64(v float64) float64 {
	if math.IsNaN(v) {
		return 0
	}
	return v
}

type recentRates struct {
	rxRateMbps          uint32
	txRateMbps          uint32
	rxRateMbpsLast15s   float32
	txRateMbpsLast15s   float32
	rxRateMbpsLast1mAvg float32
	txRateMbpsLast30s   float32
}

func mapRecentRates(rxStats *pb.WifiClient_RxStats, txStats *pb.WifiClient_TxStats) recentRates {
	rates := recentRates{
		rxRateMbps:          rxStats.GetRateMbps(),
		txRateMbps:          txStats.GetRateMbps(),
		rxRateMbpsLast15s:   rxStats.GetRateMbpsLast_15S(),
		txRateMbpsLast15s:   txStats.GetRateMbpsLast_15S(),
		rxRateMbpsLast1mAvg: rxStats.GetRateMbpsLast_30S(),
		txRateMbpsLast30s:   txStats.GetRateMbpsLast_30S(),
	}

	// Fallbacks keep recent-rate output useful when one of the near-term windows
	// is absent in a specific payload/firmware variation.
	if rates.rxRateMbpsLast15s == 0 && rates.rxRateMbps > 0 {
		rates.rxRateMbpsLast15s = float32(rates.rxRateMbps)
	}
	if rates.txRateMbpsLast15s == 0 && rates.txRateMbpsLast30s > 0 {
		rates.txRateMbpsLast15s = rates.txRateMbpsLast30s
	}
	if rates.txRateMbpsLast15s == 0 && rates.txRateMbps > 0 {
		rates.txRateMbpsLast15s = float32(rates.txRateMbps)
	}
	if rates.rxRateMbpsLast1mAvg == 0 && rates.rxRateMbpsLast15s > 0 {
		rates.rxRateMbpsLast1mAvg = rates.rxRateMbpsLast15s
	}
	if rates.txRateMbpsLast30s == 0 && rates.txRateMbpsLast15s > 0 {
		rates.txRateMbpsLast30s = rates.txRateMbpsLast15s
	}
	if rates.rxRateMbps == 0 && rates.rxRateMbpsLast15s > 0 {
		rates.rxRateMbps = uint32(rates.rxRateMbpsLast15s)
	}
	if rates.txRateMbps == 0 && rates.txRateMbpsLast15s > 0 {
		rates.txRateMbps = uint32(rates.txRateMbpsLast15s)
	}

	return rates
}

func normalizeEnum(raw string, trimPrefix string) string {
	if raw == "" {
		return ""
	}

	normalized := raw
	if trimPrefix != "" {
		normalized = strings.TrimPrefix(normalized, trimPrefix)
	}
	return strings.ToLower(normalized)
}
