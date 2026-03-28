#!/usr/bin/env bash
set -euo pipefail

#!/usr/bin/env bash
set -euo pipefail

DUMP_DIR="${1:-}"
PROJECT_ROOT="${2:-$(pwd)}"

if [[ -z "${DUMP_DIR}" ]]; then
  echo "Uso: $0 <caminho_do_dump> [caminho_do_starlink-go]" >&2
  exit 1
fi

DUMP_DIR="$(realpath "${DUMP_DIR}")"
PROJECT_ROOT="$(realpath "${PROJECT_ROOT}")"

PROTOSET="${DUMP_DIR}/dish.protoset"

if [[ ! -f "${PROTOSET}" ]]; then
  echo "Arquivo nao encontrado: ${PROTOSET}" >&2
  exit 1
fi

cd "${PROJECT_ROOT}"

if [[ ! -f "go.mod" ]]; then
  echo "go.mod nao encontrado em ${PROJECT_ROOT}" >&2
  exit 1
fi

if [[ ! -d "proto" ]]; then
  echo "pasta proto/ nao encontrada em ${PROJECT_ROOT}" >&2
  exit 1
fi

OUTDIR="${PROJECT_ROOT}/proto_sync_$(date +%Y%m%d_%H%M%S)"
mkdir -p "${OUTDIR}"/{repo_desc,protoset_desc,diff,grep}

echo "[1/8] Descrevendo schema do protoset"
grpcurl -plaintext -protoset "${PROTOSET}" describe SpaceX.API.Device.Request \
  > "${OUTDIR}/protoset_desc/Request.txt"
grpcurl -plaintext -protoset "${PROTOSET}" describe SpaceX.API.Device.Response \
  > "${OUTDIR}/protoset_desc/Response.txt"
grpcurl -plaintext -protoset "${PROTOSET}" describe SpaceX.API.Device.WifiClient \
  > "${OUTDIR}/protoset_desc/WifiClient.txt"
grpcurl -plaintext -protoset "${PROTOSET}" describe SpaceX.API.Device.WifiGetStatusResponse \
  > "${OUTDIR}/protoset_desc/WifiGetStatusResponse.txt"
grpcurl -plaintext -protoset "${PROTOSET}" describe SpaceX.API.Device.DhcpLease \
  > "${OUTDIR}/protoset_desc/DhcpLease.txt"

grpcurl -plaintext -protoset "${PROTOSET}" describe SpaceX.API.Device.NetworkInterface \
  > "${OUTDIR}/protoset_desc/NetworkInterface.txt" 2>/dev/null || true
grpcurl -plaintext -protoset "${PROTOSET}" describe SpaceX.API.Device.RadioStats \
  > "${OUTDIR}/protoset_desc/RadioStats.txt" 2>/dev/null || true

echo "[2/8] Descrevendo schema atual do repo"
grpcurl -plaintext -import-path proto \
  -proto proto/spacex/api/device/common.proto \
  -proto proto/spacex/api/device/device.proto \
  -proto proto/spacex/api/device/wifi.proto \
  describe SpaceX.API.Device.Request \
  > "${OUTDIR}/repo_desc/Request.txt"

grpcurl -plaintext -import-path proto \
  -proto proto/spacex/api/device/common.proto \
  -proto proto/spacex/api/device/device.proto \
  -proto proto/spacex/api/device/wifi.proto \
  describe SpaceX.API.Device.Response \
  > "${OUTDIR}/repo_desc/Response.txt"

grpcurl -plaintext -import-path proto \
  -proto proto/spacex/api/device/wifi.proto \
  describe SpaceX.API.Device.WifiClient \
  > "${OUTDIR}/repo_desc/WifiClient.txt"

grpcurl -plaintext -import-path proto \
  -proto proto/spacex/api/device/wifi.proto \
  describe SpaceX.API.Device.WifiGetStatusResponse \
  > "${OUTDIR}/repo_desc/WifiGetStatusResponse.txt"

grpcurl -plaintext -import-path proto \
  -proto proto/spacex/api/device/common.proto \
  describe SpaceX.API.Device.DhcpLease \
  > "${OUTDIR}/repo_desc/DhcpLease.txt" 2>"${OUTDIR}/repo_desc/DhcpLease.err" || true

grpcurl -plaintext -import-path proto \
  -proto proto/spacex/api/device/common.proto \
  -proto proto/spacex/api/device/device.proto \
  describe SpaceX.API.Device.NetworkInterface \
  > "${OUTDIR}/repo_desc/NetworkInterface.txt" 2>/dev/null || true

grpcurl -plaintext -import-path proto \
  -proto proto/spacex/api/device/wifi.proto \
  describe SpaceX.API.Device.RadioStats \
  > "${OUTDIR}/repo_desc/RadioStats.txt" 2>/dev/null || true

echo "[3/8] Gerando diffs"
for f in Request Response WifiClient WifiGetStatusResponse DhcpLease NetworkInterface RadioStats; do
  if [[ -f "${OUTDIR}/repo_desc/${f}.txt" && -f "${OUTDIR}/protoset_desc/${f}.txt" ]]; then
    diff -u "${OUTDIR}/repo_desc/${f}.txt" "${OUTDIR}/protoset_desc/${f}.txt" \
      > "${OUTDIR}/diff/${f}.diff" || true
  fi
done

echo "[4/8] Procurando campos importantes no repo atual"
grep -RniE \
  "seconds_until_dhcp_lease_expires|dhcp_lease_found|captive_client_id|upload_mb|download_mb|rx_stats_valid|tx_stats_valid|rate_mbps_last_1m_avg|throughput_mbps_last_1m_avg|captive_state|sandbox_state|alerts" \
  proto > "${OUTDIR}/grep/important_fields_in_repo.txt" || true

echo "[5/8] Procurando campos importantes nos descritores do protoset"
grep -RniE \
  "seconds_until_dhcp_lease_expires|dhcp_lease_found|captive_client_id|upload_mb|download_mb|rx_stats_valid|tx_stats_valid|rate_mbps_last_1m_avg|throughput_mbps_last_1m_avg|captive_state|sandbox_state|alerts" \
  "${OUTDIR}/protoset_desc" > "${OUTDIR}/grep/important_fields_in_protoset.txt" || true

echo "[6/8] Resumo inicial"
cat > "${OUTDIR}/SUMMARY.txt" <<EOF
Projeto: ${PROJECT_ROOT}
Dump: ${DUMP_DIR}
Protoset: ${PROTOSET}

Arquivos:
- repo_desc/
- protoset_desc/
- diff/
- grep/

Proximo passo:
1. abrir os .diff em diff/
2. atualizar manualmente os .proto em proto/spacex/api/device/
3. rodar novamente este script com --regen, ou executar:
   find proto/gen -name "*.pb.go" -type f -delete
   buf generate
   go build ./...
   go test ./...
EOF

echo "[7/8] Mostrando arquivos diff nao vazios"
find "${OUTDIR}/diff" -type f -size +0c | sort || true

echo "[8/8] Finalizado"
echo "Saida em: ${OUTDIR}"