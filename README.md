# starlink-adapter

Repositório reduzido para um SDK Go minimalista em `starlink-go/` com suporte **apenas** a acesso local via gRPC direto (porta `9200`).

## Removido

- `libs/dart/`, `libs/python/`, `libs/swift/` e demais artefatos multi-language.
- Qualquer lógica de cookies/browser/keychain/autenticação web.
- Qualquer integração grpc-web ou HTTP REST remoto.
- Features fora do core mínimo de monitoramento/operações básicas.

## Mantido

- `starlink-go/proto/` e `starlink-go/proto/gen/` para camada protobuf.
- `starlink-go/client/` para interface de domínio sem expor pb.go, com operações:
  - `GetStatus`
  - `GetStats`
  - `GetLocation`
  - `Reboot`
  - `GetConnectedClients`
- `starlink-go/internal/transport/localgrpc/` para conexão gRPC local direta.

## GetConnectedClients (uso interno)

`GetConnectedClients` retorna a visão de clientes conectados exposta pelo schema `wifi_get_clients` da Starlink (ex.: MAC/IP/interface e metadados relacionados).

Quando faz sentido usar:

- inventário técnico de clientes conectados;
- troubleshooting de conectividade;
- snapshot periódico de clientes por unidade;
- diagnóstico operacional em NOC.

Limites atuais:

- não expõe serial de clientes (o schema atual não fornece esse campo);
- representa somente o que o `wifi_get_clients` disponibiliza no momento da coleta.

Observação operacional:

- em ambientes com muitas unidades/filiais, a utilidade cresce para visão consolidada;
- a coleta deve ser moderada em larga escala (evitar polling agressivo contínuo).
