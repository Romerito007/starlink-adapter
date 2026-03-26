# starlink-adapter

Repositório reduzido para um SDK Go minimalista em `starlink-go/` com suporte **apenas** a acesso local via gRPC direto (porta `9200`).

## Removido

- `libs/dart/`, `libs/python/`, `libs/swift/` e demais artefatos multi-language.
- Qualquer lógica de cookies/browser/keychain/autenticação web.
- Qualquer integração grpc-web ou HTTP REST remoto.
- Features fora do core mínimo de monitoramento/operações básicas.

## Mantido

- `starlink-go/proto/` e `starlink-go/proto/gen/` para camada protobuf.
- `starlink-go/client/` para interface de domínio (`GetStatus`, `GetStats`, `GetLocation`, `Reboot`) sem expor pb.go.
- `starlink-go/internal/transport/localgrpc/` para conexão gRPC local direta.
