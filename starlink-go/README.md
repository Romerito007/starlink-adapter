# starlink-go

Módulo Go: `github.com/Romerito007/starlink-adapter/starlink-go`

## 1) O que é este projeto

`starlink-go` é um **adapter Go privado** para acesso técnico ao terminal Starlink via gRPC local (porta 9200), normalmente em rede local ou por malha privada (VPN/roteamento).

Este projeto **não é uma plataforma completa**: ele entrega um cliente enxuto para monitoramento e operações básicas do terminal.

Operações atualmente suportadas:

- `GetStatus(ctx context.Context) (*Status, error)`
- `GetStatusDetailed(ctx context.Context) (*StatusDetailed, error)`
- `GetStats(ctx context.Context) (*Stats, error)`
- `GetLocation(ctx context.Context) (*Location, error)`
- `GetConnectedClients(ctx context.Context) ([]ClientDevice, error)`
- `GetDhcpLeases(ctx context.Context) ([]DhcpLease, error)`
- `GetWifiConfig(ctx context.Context) (*WifiConfigSnapshot, error)`
- `GetNetworkInterfaces(ctx context.Context) ([]NetworkInterfaceSnapshot, error)`
- `GetRadioStats(ctx context.Context) ([]RadioStat, error)`
- `GetEventLogSummary(ctx context.Context) (*EventLogSummary, error)`
- `Reboot(ctx context.Context) error`
- `Close() error`

### GetConnectedClients (wifi_get_status.clients com fallback)

O adapter consulta `get_status` e, quando o endpoint responde com `wifi_get_status`, usa `wifi_get_status.clients` como fonte primária para retornar uma lista normalizada de clientes conectados (`[]ClientDevice`).

Para compatibilidade com endpoints legados, existe fallback interno para `wifi_get_clients`.

Campos atualmente mapeados:

- `MacAddress`
- `IpAddress`
- `Interface` (string normalizada legível, ex.: `eth`, `rf_2ghz`, `rf_5ghz`, `rf_5ghz_high`)
- `SignalStrength`
- `AssociatedTimeSeconds`
- `Name`
- `GivenName`
- `Domain`
- `Ipv6Addresses`
- `DhcpLeaseActive`
- `DhcpLeaseRenewed`
- `ChannelWidth`
- `Snr`
- `Mode`
- `Blocked`
- `Role`
- `InterfaceName`
- `NoDataIdleSeconds`
- `UpstreamMacAddress`
- `HopsFromController`
- `ClientID`
- `RxRateMbps`
- `TxRateMbps`
- `RxRateMbpsLast15s`
- `TxRateMbpsLast15s`
- `TxRateMbpsLast30s`
- `RxBytes`
- `TxBytes`
- `RxNss`
- `TxNss`
- `RxMcs`
- `TxMcs`
- `RxBandwidth`
- `TxBandwidth`
- `RxGuardNs`
- `TxGuardNs`
- `RxPhyMode`
- `TxPhyMode`
- `RxStatsValid`
- `TxStatsValid`

Campos operacionais adicionais solicitados, mas ainda indisponíveis no protobuf atual, permanecem com default zero/false na model pública:

- `DhcpLeaseFound`
- `SecondsUntilDhcpLeaseExpires`
- `CaptiveClientID`
- `UploadMb`
- `DownloadMb`
- `RxRateMbpsLast1mAvg`

Limitação importante: o schema protobuf atual **não expõe serial do cliente**.

### GetStatusDetailed (get_status detalhado)

`GetStatusDetailed` mantém `GetStatus` enxuto e expõe um snapshot operacional ampliado a partir de `get_status`, sem incluir blobs/config/clients.

Campos do model `StatusDetailed`:

- `DeviceID`
- `HardwareVersion`
- `SoftwareVersion`
- `UptimeSeconds`
- `Ipv4WanAddress`
- `Ipv6WanAddresses`
- `PingLatencyMs`
- `PingDropRate`
- `PingDropRate5m`
- `DishPingLatencyMs`
- `DishPingDropRate`
- `DishPingDropRate5m`
- `PopPingLatencyMs`
- `PopPingDropRate`
- `PopPingDropRate5m`
- `PopIpv6PingLatencyMs`
- `PopIpv6PingDropRate`
- `PopIpv6PingDropRate5m`
- `SecsSinceLastPublicIpv4Change`
- `DishID`
- `UtcNs`
- `DishDisablementCode`
- `CalibrationPartitionsState`
- `SetupRequirementState`
- `SoftwareUpdateState`
- `SoftwareUpdateRunningVersion`
- `SoftwareUpdateSecondsSinceGetTargetVersions`
- `PoeState`
- `PoePower`
- `PoeVin`

Observação: alguns campos acima podem permanecer com default (`0`/`""`) quando não estiverem disponíveis no protobuf atual do caminho `wifi_get_status`.

### GetDhcpLeases (wifi_get_status.dhcp_servers[].leases[])

O adapter consulta `get_status` e extrai leases DHCP a partir de `wifi_get_status.dhcp_servers[].leases[]`, retornando `[]DhcpLease` sem expor tipos protobuf.

Campos do model `DhcpLease`:

- `IpAddress`
- `MacAddress`
- `Hostname`
- `ExpiresTime`
- `Active`
- `ClientID`
- `Domain` (herdado de `dhcp_servers[].domain`)

A saída é estável e ordenada por `domain + ip_address + mac_address`.

### GetWifiConfig (wifi_get_config)

O adapter consulta `wifi_get_config` e retorna um snapshot público/normalizado de configuração Wi-Fi/LAN sem expor protobuf e sem expor campos sensíveis (como senha).

Modelos públicos:

- `WifiConfigSnapshot`
  - `CountryCode`
  - `SetupComplete`
  - `MacWan`
  - `MacLan`
  - `BootCount`
  - `Incarnation`
  - `WanHostDscpMark`
  - `Networks []WifiNetwork`
- `WifiNetwork`
  - `Ipv4`
  - `Domain`
  - `Dhcpv4Start`
  - `Dhcpv4End` (mantido no model público; no protobuf atual fica `0`)
  - `Dhcpv4LeaseDurationSeconds`
  - `Vlan`
  - `BasicServiceSets []WifiBasicServiceSet`
- `WifiBasicServiceSet`
  - `Bssid`
  - `Ssid`
  - `Band` (string legível normalizada)
  - `InterfaceName`

A saída de redes e BSS também é estável e ordenada.

### GetNetworkInterfaces (get_network_interfaces)

O adapter consulta `get_network_interfaces` e retorna um inventário normalizado de interfaces de rede do roteador (`[]NetworkInterfaceSnapshot`), sem expor protobuf.

Campos públicos:

- `NetworkInterfaceSnapshot`
  - `Name`
  - `Up`
  - `MacAddress`
  - `Ipv4Addresses`
  - `Ipv6Addresses`
  - `RxStats`
  - `TxStats`
  - `Ethernet`
  - `Wifi`
  - `Bridge`
- `InterfaceTrafficStats`
  - `Bytes`
  - `Packets`
- `InterfaceEthernetInfo`
  - `LinkDetected`
  - `SpeedMbps`
  - `AutonegotiationOn`
  - `Duplex` (string legível normalizada)
- `InterfaceWifiInfo`
  - `Channel`
  - `LinkQuality`
- `InterfaceBridgeInfo`
  - `MemberNames`

Subestruturas (`Ethernet`, `Wifi`, `Bridge`) são opcionais e nil-safe. A saída é estável e ordenada por `name`.

### GetRadioStats (get_radio_stats)

O adapter consulta `get_radio_stats` e retorna saúde básica por banda em `[]RadioStat`.

Modelos públicos:

- `RadioStat`
  - `Band`
  - `RxStats`
  - `TxStats`
  - `ThermalStatus`
  - `AntennaStatus`
- `RadioTrafficStats`
  - `Packets`
  - `FrameErrors`
- `RadioThermalStatus`
  - `Temp`
  - `DutyCycle`
- `RadioAntennaStatus`
  - `Rssi1`
  - `Rssi2`
  - `Rssi3`
  - `Rssi4`

Estratégia para NaN:

- Campos float com `NaN` vindos do payload (`thermal_status.temp2` e `antenna_status.rssi*`) são normalizados para `0` antes de expor no model público.
- Motivo: facilitar serialização/consumo downstream com comportamento determinístico.

A saída é estável e ordenada por `band`.

### GetEventLogSummary (get_history.event_log)

`GetEventLogSummary` expõe um snapshot leve de event log operacional via `get_history`, sem retornar a série histórica bruta completa.

Modelos públicos:

- `EventLogSummary`
  - `StartTimestampNs`
  - `CurrentTimestampNs`
  - `Events []EventLogEvent`
- `EventLogEvent`
  - `Severity`
  - `Reason`
  - `StartTimestampNs`
  - `DurationNs`

Observação importante:

- No protobuf atualmente versionado no submódulo, `WifiGetHistoryResponse` ainda não expõe `event_log` por getters tipados.
- Por isso, o método já retorna a estrutura pública final com nil-safety, porém com defaults (`timestamps=0`, `events=[]`) até a disponibilidade tipada desses campos.

## 2) Como a conectividade funciona

A Starlink expõe um endpoint gRPC local (tipicamente `192.168.100.1:9200` no domínio local do terminal).

Para este adapter funcionar, o serviço que executa o cliente precisa ter **alcance de rede TCP real** até o `host:port` da Starlink. Isso normalmente ocorre por:

- mesma LAN da Starlink; ou
- VPN (site-to-site, hub-and-spoke etc.); ou
- roteamento entre redes internas até a rede onde está a Starlink.

Esta lib **não depende de API web remota** e **não usa cookie/browser**.

## 3) Cenário recomendado com MikroTik (exemplo principal)

Exemplo típico em produção:

1. A unidade/filial do cliente possui um roteador (ex.: MikroTik) conectado à rede local da Starlink.
2. Esse roteador estabelece VPN para a infraestrutura central.
3. A infraestrutura central passa a alcançar o IP local da Starlink (ou um endpoint redirecionado) via malha privada.
4. O serviço central usa este adapter apontando para o `host:port` alcançável.

Quando usar port-forward/NAT interno:

- quando o serviço não enxerga diretamente o IP local da Starlink;
- quando a topologia exige um endpoint intermediário no roteador da unidade.

Ponto principal: o adapter só precisa que o `host:port` final esteja alcançável pela rede privada.

Sobre IP fixo público:

- **não é obrigatório em toda topologia**;
- com VPN bem montada, o acesso pode ser estável sem IP público fixo;
- IP fixo público faz sentido quando ajuda em previsibilidade de túneis, simplificação operacional e troubleshooting.

## 4) Requisitos de rede

Para operação estável:

- reachability TCP até `host:port` da Starlink;
- rotas corretas no roteador/firewall;
- ACL/firewall permitindo tráfego necessário entre serviço e terminal;
- latência e estabilidade de rede compatíveis com chamadas gRPC.

> Observação: `192.168.100.1:9200` é o padrão comum local, mas o endpoint efetivo pode variar conforme a topologia (VPN/NAT/roteamento).

## 5) Exemplo de configuração lógica (arquitetura)

```text
[Serviço central / NOC]
          |
          | (VPN privada)
          v
[Roteador da unidade - ex. MikroTik]
          |
          | (LAN local da unidade)
          v
[Starlink endpoint gRPC: 192.168.100.1:9200]
```

Fluxo operacional:

- o serviço central envia chamadas gRPC para o `host:port` configurado;
- a malha de conectividade (VPN + roteamento + firewall) entrega o tráfego até a rede da Starlink;
- o adapter executa chamadas de status/estatísticas/localização/reboot.

## 6) Como configurar o client

Construtor público:

```go
func NewClient(ctx context.Context, cfg Config) (StarlinkClient, error)
```

Config disponível:

- `Host string`
- `Port int`
- `Timeout time.Duration`
- `Logger *slog.Logger` (opcional)

Defaults técnicos aplicados quando não informados:

- `Host`: `192.168.100.1`
- `Port`: `9200`
- `Timeout`: `5s`

`Host`/`Port` devem apontar para o endpoint **realmente alcançável** pelo serviço (LAN/VPN/roteamento privado).

## 7) Quando este adapter faz sentido

Este adapter é útil quando você já possui conectividade de rede entre um sistema central e muitos terminais, por exemplo:

- operações de NOC;
- monitoramento centralizado de filiais/unidades;
- reboot controlado em operação;
- coleta periódica de saúde/métricas.

Escalas típicas de uso: ambientes com 100, 300 ou 1000+ unidades, desde que exista malha de conectividade bem definida.

O ganho principal é padronizar acesso técnico ao terminal dentro de um ecossistema maior de observabilidade/orquestração.

## 8) Limitações

Este adapter **não**:

- faz descoberta automática de rede;
- cria ou gerencia VPN;
- abre firewall/rota automaticamente;
- resolve NAT/topologia por conta própria;
- faz autenticação remota web de conta Starlink.

Ele depende de reachability real até o terminal e deve ser combinado com um sistema maior de monitoramento/provisioning.

## 9) Exemplo mínimo de uso

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Romerito007/starlink-adapter/starlink-go/client"
)

func main() {
	cli, err := client.NewClient(context.Background(), client.Config{
		Host:    "192.168.100.1", // ou endpoint alcançável via VPN/roteamento
		Port:    9200,
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
		fmt.Printf("mac=%s ip=%s iface=%s signal=%.1f rx=%d tx=%d blocked=%t\n",
			c.MacAddress, c.IpAddress, c.Interface, c.SignalStrength, c.RxRateMbps, c.TxRateMbps, c.Blocked)
	}
}
```

## 10) Observações operacionais

- Ajuste `Timeout` por perfil de rede (latência/jitter entre central e unidade).
- O client já aplica retry/backoff simples para falhas transitórias, mas isso não substitui uma malha de rede estável.
- Feche o client com `Close()` ao encerrar worker/job/processo.
- Use `GetConnectedClients` com moderação em polling; prefira snapshots periódicos para inventário e diagnóstico, não loop agressivo contínuo.
- Evite polling agressivo em larga escala; prefira agendamento controlado, filas e workers.
- Em ambientes grandes, distribua coleta por lotes para reduzir picos de carga e facilitar troubleshooting.
