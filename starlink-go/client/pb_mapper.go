package client

import (
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
