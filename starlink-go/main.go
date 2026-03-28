package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Romerito007/starlink-adapter/starlink-go/client"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cli, err := client.NewClient(ctx, client.Config{
		Host:    "100.126.255.11", // endpoint alcançável via VPN/roteamento
		Port:    9000,             // 9200 dish, 9000 router
		Timeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	fmt.Println("=== STATUS ===")
	status, err := cli.GetStatus(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("device_id=%q\n", status.DeviceID)
	fmt.Printf("hardware_version=%q\n", status.HardwareVersion)
	fmt.Printf("software_version=%q\n", status.SoftwareVersion)
	fmt.Printf("uptime_seconds=%d\n", status.UptimeSeconds)
	fmt.Printf("uplink_throughput_bps=%.2f\n", status.UplinkThroughputBps)
	fmt.Printf("downlink_throughput_bps=%.2f\n", status.DownlinkThroughputBps)
	fmt.Printf("pop_ping_drop_rate=%.6f\n", status.PopPingDropRate)
	fmt.Printf("pop_ping_latency_ms=%.2f\n", status.PopPingLatencyMs)
	fmt.Println()

	fmt.Println("=== CONNECTED CLIENTS ===")
	clients, err := cli.GetConnectedClients(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("count=%d\n", len(clients))
	for i, c := range clients {
		fmt.Printf("client[%d]\n", i)
		fmt.Printf("  mac_address=%q\n", c.MacAddress)
		fmt.Printf("  ip_address=%q\n", c.IpAddress)
		fmt.Printf("  interface=%q\n", c.Interface)
		fmt.Printf("  interface_name=%q\n", c.InterfaceName)
		fmt.Printf("  upstream_mac_address=%q\n", c.UpstreamMacAddress)
		fmt.Printf("  role=%q\n", c.Role)
		fmt.Printf("  signal_strength=%.3f\n", c.SignalStrength)
		fmt.Printf("  snr=%.3f\n", c.Snr)
		fmt.Printf("  channel_width=%d\n", c.ChannelWidth)
		fmt.Printf("  mode=%q\n", c.Mode)
		fmt.Printf("  blocked=%t\n", c.Blocked)
		fmt.Printf("  dhcp_lease_active=%t\n", c.DhcpLeaseActive)
		fmt.Printf("  dhcp_lease_renewed=%t\n", c.DhcpLeaseRenewed)
		fmt.Printf("  dhcp_lease_found=%t\n", c.DhcpLeaseFound)
		fmt.Printf("  seconds_until_dhcp_lease_expires=%d\n", c.SecondsUntilDhcpLeaseExpires)
		fmt.Printf("  associated_time_seconds=%d\n", c.AssociatedTimeSeconds)
		fmt.Printf("  no_data_idle_seconds=%d\n", c.NoDataIdleSeconds)
		fmt.Printf("  hops_from_controller=%d\n", c.HopsFromController)
		fmt.Printf("  client_id=%d\n", c.ClientID)
		fmt.Printf("  captive_client_id=%q\n", c.CaptiveClientID)
		fmt.Printf("  upload_mb=%.3f\n", c.UploadMb)
		fmt.Printf("  download_mb=%.3f\n", c.DownloadMb)
		fmt.Printf("  rx_stats_valid=%t\n", c.RxStatsValid)
		fmt.Printf("  tx_stats_valid=%t\n", c.TxStatsValid)
		fmt.Printf("  rx_bytes=%d\n", c.RxBytes)
		fmt.Printf("  tx_bytes=%d\n", c.TxBytes)
		fmt.Printf("  rx_nss=%d\n", c.RxNss)
		fmt.Printf("  tx_nss=%d\n", c.TxNss)
		fmt.Printf("  rx_mcs=%d\n", c.RxMcs)
		fmt.Printf("  tx_mcs=%d\n", c.TxMcs)
		fmt.Printf("  rx_bandwidth=%d\n", c.RxBandwidth)
		fmt.Printf("  tx_bandwidth=%d\n", c.TxBandwidth)
		fmt.Printf("  rx_guard_ns=%d\n", c.RxGuardNs)
		fmt.Printf("  tx_guard_ns=%d\n", c.TxGuardNs)
		fmt.Printf("  rx_rate_mbps=%d\n", c.RxRateMbps)
		fmt.Printf("  tx_rate_mbps=%d\n", c.TxRateMbps)
		fmt.Printf("  rx_phy_mode=%d\n", c.RxPhyMode)
		fmt.Printf("  tx_phy_mode=%d\n", c.TxPhyMode)
		fmt.Printf("  rx_rate_mbps_last_15s=%.3f\n", c.RxRateMbpsLast15s)
		fmt.Printf("  tx_rate_mbps_last_15s=%.3f\n", c.TxRateMbpsLast15s)
		fmt.Printf("  rx_rate_mbps_last_1m_avg=%.3f\n", c.RxRateMbpsLast1mAvg)
		fmt.Printf("  tx_rate_mbps_last_30s=%.3f\n", c.TxRateMbpsLast30s)
		fmt.Printf("  name=%q\n", c.Name)
		fmt.Printf("  given_name=%q\n", c.GivenName)
		fmt.Printf("  domain=%q\n", c.Domain)
		fmt.Printf("  ipv6_addresses=%v\n", c.Ipv6Addresses)
	}
	fmt.Println()

	fmt.Println("=== DHCP LEASES ===")
	leases, err := cli.GetDhcpLeases(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("count=%d\n", len(leases))
	for i, l := range leases {
		fmt.Printf("lease[%d]\n", i)
		fmt.Printf("  ip_address=%q\n", l.IpAddress)
		fmt.Printf("  mac_address=%q\n", l.MacAddress)
		fmt.Printf("  hostname=%q\n", l.Hostname)
		fmt.Printf("  expires_time=%q\n", l.ExpiresTime)
		fmt.Printf("  active=%t\n", l.Active)
		fmt.Printf("  client_id=%d\n", l.ClientID)
		fmt.Printf("  domain=%q\n", l.Domain)
	}
	fmt.Println()

	fmt.Println("=== WIFI CONFIG ===")
	wifiCfg, err := cli.GetWifiConfig(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("country_code=%q\n", wifiCfg.CountryCode)
	fmt.Printf("setup_complete=%t\n", wifiCfg.SetupComplete)
	fmt.Printf("mac_wan=%q\n", wifiCfg.MacWan)
	fmt.Printf("mac_lan=%q\n", wifiCfg.MacLan)
	fmt.Printf("boot_count=%d\n", wifiCfg.BootCount)
	fmt.Printf("incarnation=%d\n", wifiCfg.Incarnation)
	fmt.Printf("wan_host_dscp_mark=%d\n", wifiCfg.WanHostDscpMark)
	fmt.Printf("networks_count=%d\n", len(wifiCfg.Networks))
	for i, n := range wifiCfg.Networks {
		fmt.Printf("network[%d]\n", i)
		fmt.Printf("  ipv4=%q\n", n.Ipv4)
		fmt.Printf("  domain=%q\n", n.Domain)
		fmt.Printf("  dhcpv4_start=%d\n", n.Dhcpv4Start)
		fmt.Printf("  dhcpv4_end=%d\n", n.Dhcpv4End)
		fmt.Printf("  dhcpv4_lease_duration_seconds=%d\n", n.Dhcpv4LeaseDurationSeconds)
		fmt.Printf("  vlan=%d\n", n.Vlan)
		fmt.Printf("  basic_service_sets_count=%d\n", len(n.BasicServiceSets))
		for j, b := range n.BasicServiceSets {
			fmt.Printf("  basic_service_set[%d]\n", j)
			fmt.Printf("    bssid=%q\n", b.Bssid)
			fmt.Printf("    ssid=%q\n", b.Ssid)
			fmt.Printf("    band=%q\n", b.Band)
			fmt.Printf("    interface_name=%q\n", b.InterfaceName)
		}
	}
	fmt.Println()

	fmt.Println("=== NETWORK INTERFACES ===")
	ifaces, err := cli.GetNetworkInterfaces(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("count=%d\n", len(ifaces))
	for i, ni := range ifaces {
		fmt.Printf("interface[%d]\n", i)
		fmt.Printf("  name=%q\n", ni.Name)
		fmt.Printf("  up=%t\n", ni.Up)
		fmt.Printf("  mac_address=%q\n", ni.MacAddress)
		fmt.Printf("  ipv4_addresses=%v\n", ni.Ipv4Addresses)
		fmt.Printf("  ipv6_addresses=%v\n", ni.Ipv6Addresses)
		fmt.Printf("  rx_stats.bytes=%d\n", ni.RxStats.Bytes)
		fmt.Printf("  rx_stats.packets=%d\n", ni.RxStats.Packets)
		fmt.Printf("  tx_stats.bytes=%d\n", ni.TxStats.Bytes)
		fmt.Printf("  tx_stats.packets=%d\n", ni.TxStats.Packets)
		if ni.Ethernet != nil {
			fmt.Printf("  ethernet.link_detected=%t\n", ni.Ethernet.LinkDetected)
			fmt.Printf("  ethernet.speed_mbps=%d\n", ni.Ethernet.SpeedMbps)
			fmt.Printf("  ethernet.autonegotiation_on=%t\n", ni.Ethernet.AutonegotiationOn)
			fmt.Printf("  ethernet.duplex=%q\n", ni.Ethernet.Duplex)
		}
		if ni.Wifi != nil {
			fmt.Printf("  wifi.channel=%d\n", ni.Wifi.Channel)
			fmt.Printf("  wifi.link_quality=%.6f\n", ni.Wifi.LinkQuality)
		}
		if ni.Bridge != nil {
			fmt.Printf("  bridge.member_names=%v\n", ni.Bridge.MemberNames)
		}
	}
	fmt.Println()

	fmt.Println("=== RADIO STATS ===")
	radioStats, err := cli.GetRadioStats(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("count=%d\n", len(radioStats))
	for i, r := range radioStats {
		fmt.Printf("radio_stat[%d]\n", i)
		fmt.Printf("  band=%q\n", r.Band)
		fmt.Printf("  rx_stats.packets=%d\n", r.RxStats.Packets)
		fmt.Printf("  rx_stats.frame_errors=%d\n", r.RxStats.FrameErrors)
		fmt.Printf("  tx_stats.packets=%d\n", r.TxStats.Packets)
		fmt.Printf("  tx_stats.frame_errors=%d\n", r.TxStats.FrameErrors)
		fmt.Printf("  thermal_status.temp=%.6f\n", r.ThermalStatus.Temp)
		fmt.Printf("  thermal_status.duty_cycle=%d\n", r.ThermalStatus.DutyCycle)
		fmt.Printf("  antenna_status.rssi_1=%.6f\n", r.AntennaStatus.Rssi1)
		fmt.Printf("  antenna_status.rssi_2=%.6f\n", r.AntennaStatus.Rssi2)
		fmt.Printf("  antenna_status.rssi_3=%.6f\n", r.AntennaStatus.Rssi3)
		fmt.Printf("  antenna_status.rssi_4=%.6f\n", r.AntennaStatus.Rssi4)
	}
	fmt.Println()

	fmt.Println("=== STATUS DETAILED ===")
	detailed, err := cli.GetStatusDetailed(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("device_id=%q\n", detailed.DeviceID)
	fmt.Printf("hardware_version=%q\n", detailed.HardwareVersion)
	fmt.Printf("software_version=%q\n", detailed.SoftwareVersion)
	fmt.Printf("uptime_seconds=%d\n", detailed.UptimeSeconds)
	fmt.Printf("ipv4_wan_address=%q\n", detailed.Ipv4WanAddress)
	fmt.Printf("ipv6_wan_addresses=%v\n", detailed.Ipv6WanAddresses)
	fmt.Printf("ping_latency_ms=%.6f\n", detailed.PingLatencyMs)
	fmt.Printf("ping_drop_rate=%.6f\n", detailed.PingDropRate)
	fmt.Printf("ping_drop_rate_5m=%.6f\n", detailed.PingDropRate5m)
	fmt.Printf("dish_ping_latency_ms=%.6f\n", detailed.DishPingLatencyMs)
	fmt.Printf("dish_ping_drop_rate=%.6f\n", detailed.DishPingDropRate)
	fmt.Printf("dish_ping_drop_rate_5m=%.6f\n", detailed.DishPingDropRate5m)
	fmt.Printf("pop_ping_latency_ms=%.6f\n", detailed.PopPingLatencyMs)
	fmt.Printf("pop_ping_drop_rate=%.6f\n", detailed.PopPingDropRate)
	fmt.Printf("pop_ping_drop_rate_5m=%.6f\n", detailed.PopPingDropRate5m)
	fmt.Printf("pop_ipv6_ping_latency_ms=%.6f\n", detailed.PopIpv6PingLatencyMs)
	fmt.Printf("pop_ipv6_ping_drop_rate=%.6f\n", detailed.PopIpv6PingDropRate)
	fmt.Printf("pop_ipv6_ping_drop_rate_5m=%.6f\n", detailed.PopIpv6PingDropRate5m)
	fmt.Printf("secs_since_last_public_ipv4_change=%.6f\n", detailed.SecsSinceLastPublicIpv4Change)
	fmt.Printf("dish_id=%q\n", detailed.DishID)
	fmt.Printf("utc_ns=%d\n", detailed.UtcNs)
	fmt.Printf("dish_disablement_code=%q\n", detailed.DishDisablementCode)
	fmt.Printf("calibration_partitions_state=%q\n", detailed.CalibrationPartitionsState)
	fmt.Printf("setup_requirement_state=%q\n", detailed.SetupRequirementState)
	fmt.Printf("software_update_state=%q\n", detailed.SoftwareUpdateState)
	fmt.Printf("software_update_running_version=%q\n", detailed.SoftwareUpdateRunningVersion)
	fmt.Printf("software_update_seconds_since_get_target_versions=%.6f\n", detailed.SoftwareUpdateSecondsSinceGetTargetVersions)
	fmt.Printf("poe_state=%q\n", detailed.PoeState)
	fmt.Printf("poe_power=%.6f\n", detailed.PoePower)
	fmt.Printf("poe_vin=%.6f\n", detailed.PoeVin)
	fmt.Println()

	fmt.Println("=== EVENT LOG SUMMARY ===")
	eventSummary, err := cli.GetEventLogSummary(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("start_timestamp_ns=%d\n", eventSummary.StartTimestampNs)
	fmt.Printf("current_timestamp_ns=%d\n", eventSummary.CurrentTimestampNs)
	fmt.Printf("events_count=%d\n", len(eventSummary.Events))
	for i, e := range eventSummary.Events {
		fmt.Printf("event[%d]\n", i)
		fmt.Printf("  severity=%q\n", e.Severity)
		fmt.Printf("  reason=%q\n", e.Reason)
		fmt.Printf("  start_timestamp_ns=%d\n", e.StartTimestampNs)
		fmt.Printf("  duration_ns=%d\n", e.DurationNs)
	}
}
