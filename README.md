# starlink-adapter

Repositório reduzido para um SDK Go minimalista em `starlink-go/` para comunicação gRPC local com o Starlink dish (`192.168.100.1:9200`).

## Removido

- `libs/dart/`, `libs/python/`, `libs/swift/` e demais artefatos multi-language.
- Código de grpc-web, cookies, browser e autenticação remota do cliente antigo.
- scripts/docs/build tooling legados (`scripts/`, `docs/`, `Makefile`, configurações buf antigas).

## Mantido

- `starlink-go/proto/` com os `.proto` do dispositivo.
- `starlink-go/proto/gen/` com mensagens Go protobuf usadas pelo cliente.
- `starlink-go/client/` com cliente gRPC local mínimo.
