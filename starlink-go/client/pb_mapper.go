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

	state := ""
	if deviceState != nil {
		state = strconv.FormatUint(deviceState.GetUptimeS(), 10)
	}

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

		out = append(out, ClientDevice{
			ClientID:              c.GetClientId(),
			Name:                  c.GetName(),
			GivenName:             c.GetGivenName(),
			Domain:                c.GetDomain(),
			MacAddress:            c.GetMacAddress(),
			IpAddress:             c.GetIpAddress(),
			Ipv6Addresses:         ipv6,
			UpstreamMacAddress:    c.GetUpstreamMacAddress(),
			AssociatedTimeSeconds: c.GetAssociatedTimeS(),
			SignalStrength:        c.GetSignalStrength(),
			Interface:             c.GetIface().String(),
			InterfaceName:         c.GetIfaceName(),
			Role:                  c.GetRole().String(),
			DeviceID:              c.GetDeviceId(),
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

		return out[i].ClientID < out[j].ClientID
	})

	return out
}
