#!/usr/bin/env bash
set -euo pipefail

HOST="${1:-100.126.255.11}"
PORT="${2:-9000}"
PROJECT_ROOT="${3:-./starlink-go}"
TIMESTAMP="$(date +%Y%m%d_%H%M%S)"

ROOT_DIR="$(realpath "$(pwd)")"
PROJECT_ROOT="$(realpath "${PROJECT_ROOT}")"
DUMP_DIR="${PROJECT_ROOT}/starlink_dump_${TIMESTAMP}"

ROOT_DIR="$(pwd)"

echo "[1/7] Validando estrutura do projeto"
if [[ ! -d "${PROJECT_ROOT}" ]]; then
  echo "Projeto nao encontrado: ${PROJECT_ROOT}" >&2
  exit 1
fi

if [[ ! -f "${ROOT_DIR}/dump_starlink_proto.sh" ]]; then
  echo "Script dump_starlink_proto.sh nao encontrado na raiz do repo" >&2
  exit 1
fi

if [[ ! -f "${ROOT_DIR}/sync_starlink_proto.sh" ]]; then
  echo "Script sync_starlink_proto.sh nao encontrado na raiz do repo" >&2
  exit 1
fi

echo "[2/7] Rodando dump do endpoint Starlink ${HOST}:${PORT}"
"${ROOT_DIR}/dump_starlink_proto.sh" "${HOST}" "${PORT}" "${DUMP_DIR}"

echo "[3/7] Rodando comparacao/sync do protoset contra o projeto"
"${ROOT_DIR}/sync_starlink_proto.sh" "${DUMP_DIR}" "${PROJECT_ROOT}"

echo "[4/7] Limpando gerados antigos"
cd "${PROJECT_ROOT}"
find proto/gen -name "*.pb.go" -type f -delete

echo "[5/7] Regenerando protobuf com buf"
buf generate

echo "[6/7] Buildando projeto"
go build ./...

echo "[7/7] Rodando testes"
go test ./...

echo
echo "Concluido com sucesso."
echo "Dump gerado em: ${DUMP_DIR}"
echo "Projeto validado em: ${PROJECT_ROOT}"