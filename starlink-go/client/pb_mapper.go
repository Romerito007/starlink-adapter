package client

import pb "github.com/Romerito007/starlink-adapter/starlink-go/proto/gen/spacex/api/device"

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
		State:                 deviceState.GetUptimeState().String(),
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
	return &Location{
		Latitude:  lla.GetLat(),
		Longitude: lla.GetLon(),
		Altitude:  lla.GetAlt(),
		SigmaM:    in.GetSigmaM(),
		Source:    in.GetSource().String(),
	}
}
