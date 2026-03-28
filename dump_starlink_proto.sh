#!/usr/bin/env bash
set -euo pipefail

HOST="${1:-192.168.100.1}"
PORT="${2:-9000}"
OUTDIR="${3:-starlink_dump_$(date +%Y%m%d_%H%M%S)}"

if ! command -v grpcurl >/dev/null 2>&1; then
  echo "grpcurl nao encontrado no PATH" >&2
  exit 1
fi

mkdir -p "${OUTDIR}"/{describes,json}

echo "[1/6] Testando TCP em ${HOST}:${PORT}"
if command -v nc >/dev/null 2>&1; then
  nc -vz "${HOST}" "${PORT}" || {
    echo "Falha ao conectar em ${HOST}:${PORT}" >&2
    exit 1
  }
else
  echo "nc nao encontrado; pulando teste TCP"
fi

echo "[2/6] Extraindo protoset"
grpcurl -plaintext -protoset-out "${OUTDIR}/dish.protoset" "${HOST}:${PORT}" describe SpaceX.API.Device.Device

echo "[3/6] Listando services"
grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" list \
  > "${OUTDIR}/services.txt"

echo "[4/6] Salvando describe de mensagens principais"
grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" describe SpaceX.API.Device.Device \
  > "${OUTDIR}/describes/Device.service.txt"
grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" describe SpaceX.API.Device.Request \
  > "${OUTDIR}/describes/Request.txt"
grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" describe SpaceX.API.Device.Response \
  > "${OUTDIR}/describes/Response.txt"
grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" describe SpaceX.API.Device.WifiClient \
  > "${OUTDIR}/describes/WifiClient.txt"
grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" describe SpaceX.API.Device.WifiGetStatusResponse \
  > "${OUTDIR}/describes/WifiGetStatusResponse.txt"
grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" describe SpaceX.API.Device.DhcpLease \
  > "${OUTDIR}/describes/DhcpLease.txt"

# Alguns describes podem nao existir dependendo do endpoint/schema; ignorar falha
grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" describe SpaceX.API.Device.NetworkInterface \
  > "${OUTDIR}/describes/NetworkInterface.txt" 2>/dev/null || true
grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" describe SpaceX.API.Device.RadioStats \
  > "${OUTDIR}/describes/RadioStats.txt" 2>/dev/null || true

echo "[5/6] Salvando respostas JSON uteis"
grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" \
  -d '{"get_device_info":{}}' \
  "${HOST}:${PORT}" SpaceX.API.Device.Device/Handle \
  > "${OUTDIR}/json/get_device_info.json" 2>"${OUTDIR}/json/get_device_info.err" || true

grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" \
  -d '{"get_status":{}}' \
  "${HOST}:${PORT}" SpaceX.API.Device.Device/Handle \
  > "${OUTDIR}/json/get_status.json" 2>"${OUTDIR}/json/get_status.err" || true

grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" \
  -d '{"get_diagnostics":{}}' \
  "${HOST}:${PORT}" SpaceX.API.Device.Device/Handle \
  > "${OUTDIR}/json/get_diagnostics.json" 2>"${OUTDIR}/json/get_diagnostics.err" || true

grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" \
  -d '{"get_network_interfaces":{}}' \
  "${HOST}:${PORT}" SpaceX.API.Device.Device/Handle \
  > "${OUTDIR}/json/get_network_interfaces.json" 2>"${OUTDIR}/json/get_network_interfaces.err" || true

grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" \
  -d '{"get_radio_stats":{}}' \
  "${HOST}:${PORT}" SpaceX.API.Device.Device/Handle \
  > "${OUTDIR}/json/get_radio_stats.json" 2>"${OUTDIR}/json/get_radio_stats.err" || true

grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" \
  -d '{"wifi_get_config":{}}' \
  "${HOST}:${PORT}" SpaceX.API.Device.Device/Handle \
  > "${OUTDIR}/json/wifi_get_config.json" 2>"${OUTDIR}/json/wifi_get_config.err" || true

grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" \
  -d '{"wifi_get_clients":{}}' \
  "${HOST}:${PORT}" SpaceX.API.Device.Device/Handle \
  > "${OUTDIR}/json/wifi_get_clients.json" 2>"${OUTDIR}/json/wifi_get_clients.err" || true

grpcurl -plaintext -protoset "${OUTDIR}/dish.protoset" \
  -d '{"get_history":{}}' \
  "${HOST}:${PORT}" SpaceX.API.Device.Device/Handle \
  > "${OUTDIR}/json/get_history.json" 2>"${OUTDIR}/json/get_history.err" || true

echo "[6/6] Gerando indice"
cat > "${OUTDIR}/README.txt" <<EOF
Host: ${HOST}
Port: ${PORT}

Arquivos principais:
- dish.protoset
- services.txt
- describes/
- json/

Se quiser copiar para o projeto:
cp -r "${OUTDIR}" /caminho/do/seu/projeto/
EOF

echo
echo "Dump concluido em: ${OUTDIR}"
find "${OUTDIR}" -maxdepth 2 -type f | sort