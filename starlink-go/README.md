# starlink-go

SDK Go minimalista para acesso local à API gRPC do Starlink dish (`192.168.100.1:9200`).

## Estrutura

- `client/`: cliente mínimo para conectar e chamar `Device.Handle`.
- `proto/`: arquivos `.proto` originais.
- `proto/gen/`: código Go gerado necessário para serialização protobuf.

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
fmt.Println(status.DeviceInfo.Id)
```
