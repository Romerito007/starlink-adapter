package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Romerito007/starlink-adapter/starlink-go/client"
)

func main() {
	cli, err := client.NewClient(context.Background(), client.Config{
		Host:    "45.172.144.97", // ou endpoint alcançável via VPN/roteamento 92000 9000
		Port:    9000,
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

	fmt.Printf("device_id=%s uptime=%d\n", status.DeviceID, status.UptimeSeconds)

	clients, err := cli.GetConnectedClients(context.Background())
	if err != nil {
		panic(err)
	}

	for _, c := range clients {
		fmt.Printf("client_id=%d mac=%s ip=%s iface=%s signal=%.1f\n",
			c.Name, c.MacAddress, c.IpAddress, c.Interface, c.SignalStrength)
	}
}
