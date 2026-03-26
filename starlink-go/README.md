# starlink-go

SDK Go minimalista e **somente local** para acesso à API gRPC do Starlink dish (`192.168.100.1:9200`, inclusive via VPN para a LAN).

## Escopo do cliente

A API pública expõe apenas operações básicas:

- `GetStatus`
- `GetStats`
- `GetLocation`
- `Reboot`

Sem exposição direta de structs protobuf para o consumidor.

## Garantias deste SDK

- Sem cookies, sem leitura de browser, sem keychain.
- Sem autenticação web/token/refresh.
- Sem grpc-web e sem REST remoto.
- Transporte único: gRPC direto (HTTP/2) para `Device.Handle`.

## Estrutura

- `client/`: interface `StarlinkClient`, modelos de domínio e implementação `grpcClient`.
- `internal/transport/localgrpc/`: implementação de transporte gRPC local.
- `proto/`: arquivos `.proto` do dispositivo.
- `proto/gen/`: código Go protobuf para camada de transporte/protocolo.

## Exemplo rápido

```go
ctx := context.Background()
cli, err := client.Dial(ctx, "")
if err != nil {
    panic(err)
}
defer cli.Close()

status, err := cli.GetStatus(ctx)
if err != nil {
    panic(err)
}
fmt.Println(status.DeviceID)
```
