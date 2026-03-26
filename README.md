# starlink-adapter

<<<<<<< codex/refactor-repository-for-minimal-go-grpc-client-6spq1y
Repositório reduzido para um SDK Go minimalista em `starlink-go/` com suporte **apenas** a acesso local via gRPC direto (porta `9200`).
=======
Repositório reduzido para um SDK Go minimalista em `starlink-go/` para comunicação gRPC local com o Starlink dish (`192.168.100.1:9200`).
>>>>>>> main

## Removido

- `libs/dart/`, `libs/python/`, `libs/swift/` e demais artefatos multi-language.
<<<<<<< codex/refactor-repository-for-minimal-go-grpc-client-6spq1y
- Qualquer lógica de cookies/browser/keychain/autenticação web.
- Qualquer integração grpc-web ou HTTP REST remoto.
=======
- Código de grpc-web, cookies, browser e autenticação remota do cliente antigo.
>>>>>>> main
- scripts/docs/build tooling legados (`scripts/`, `docs/`, `Makefile`, configurações buf antigas).

## Mantido

- `starlink-go/proto/` com os `.proto` do dispositivo.
- `starlink-go/proto/gen/` com mensagens Go protobuf usadas pelo cliente.
<<<<<<< codex/refactor-repository-for-minimal-go-grpc-client-6spq1y
- `starlink-go/client/` + `starlink-go/internal/transport/localgrpc/` para conexão gRPC local direta.
=======
- `starlink-go/client/` com cliente gRPC local mínimo.
>>>>>>> main
