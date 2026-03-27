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

	status, err := cli.GetStatus(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"device_id=%s hardware=%s software=%s uptime=%d uplink_bps=%.2f downlink_bps=%.2f pop_drop=%.4f pop_latency_ms=%.2f\n",
		status.DeviceID,
		status.HardwareVersion,
		status.SoftwareVersion,
		status.UptimeSeconds,
		status.UplinkThroughputBps,
		status.DownlinkThroughputBps,
		status.PopPingDropRate,
		status.PopPingLatencyMs,
	)

	clients, err := cli.GetConnectedClients(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("connected_clients=%d\n", len(clients))

	for i, c := range clients {
		fmt.Printf("client[%d]\n", i)
		fmt.Printf("  name=%q\n", c.Name)
		fmt.Printf("  given_name=%q\n", c.GivenName)
		fmt.Printf("  domain=%q\n", c.Domain)
		fmt.Printf("  mac_address=%s\n", c.MacAddress)
		fmt.Printf("  ip_address=%s\n", c.IpAddress)
		fmt.Printf("  ipv6_addresses=%v\n", c.Ipv6Addresses)
		fmt.Printf("  interface=%s\n", c.Interface)
		fmt.Printf("  signal_strength=%.1f\n", c.SignalStrength)
		fmt.Printf("  associated_time_seconds=%d\n", c.AssociatedTimeSeconds)

		// Campos extended — só vão compilar se já existirem no teu ClientDevice atual.
		fmt.Printf("  dhcp_lease_active=%t\n", c.DhcpLeaseActive)
		fmt.Printf("  dhcp_lease_renewed=%t\n", c.DhcpLeaseRenewed)
		fmt.Printf("  channel_width=%d\n", c.ChannelWidth)
		fmt.Printf("  snr=%.2f\n", c.Snr)
		fmt.Printf("  mode=%q\n", c.Mode)
		fmt.Printf("  blocked=%t\n", c.Blocked)
		fmt.Printf("  role=%q\n", c.Role)
		fmt.Printf("  interface_name=%q\n", c.InterfaceName)
		fmt.Printf("  no_data_idle_seconds=%d\n", c.NoDataIdleSeconds)
		fmt.Printf("  rx_rate_mbps=%d\n", c.RxRateMbps)
		fmt.Printf("  tx_rate_mbps=%d\n", c.TxRateMbps)
		fmt.Printf("  rx_rate_mbps_last_15s=%.2f\n", c.RxRateMbpsLast15s)
		fmt.Printf("  tx_rate_mbps_last_15s=%.2f\n", c.TxRateMbpsLast15s)
	}
}
