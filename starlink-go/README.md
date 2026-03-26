# starlink-go

Módulo Go: `github.com/Romerito007/starlink-adapter/starlink-go`

Cliente Go para acesso **local gRPC** ao terminal Starlink (Dishy).

## Objetivo do pacote

Fornecer um cliente Go com API pública enxuta para operações básicas de leitura/controle da Starlink via endpoint local `host:port` (tipicamente LAN, ou VPN com alcance da LAN).

## Configuração

O cliente usa:

- `Host string`
- `Port int`
- `Timeout time.Duration`

Defaults aplicados por `NewClient` quando não informados:

- `Host`: `192.168.100.1`
- `Port`: `9200`
- `Timeout`: `5s`

## Construtor público

```go
func NewClient(cfg Config) (StarlinkClient, error)
```

## Operações suportadas

A interface pública `StarlinkClient` expõe somente:

- `GetStatus(ctx context.Context) (*Status, error)`
- `GetStats(ctx context.Context) (*Stats, error)`
- `GetLocation(ctx context.Context) (*Location, error)`
- `Reboot(ctx context.Context) error`

## Exemplo mínimo

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Romerito007/starlink-adapter/starlink-go/client"
)

func main() {
	cli, err := client.NewClient(client.Config{
		Host:    "192.168.100.1",
		Port:    9200,
		Timeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	status, err := cli.GetStatus(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(status.DeviceID)
}
```

## Observação de rede

O transporte suportado é apenas gRPC local (`internal/transport/localgrpc`), com acesso esperado por rede local/VPN. Não há API remota web neste pacote.
