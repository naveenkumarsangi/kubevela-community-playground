#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
VELA_BIN="${VELA_BIN:-vela}"
GENERATED_DIR="${ROOT_DIR}/generated"

for command in go "${VELA_BIN}"; do
  command -v "${command}" >/dev/null || {
    echo "missing required command: ${command}" >&2
    exit 1
  }
done

cd "${ROOT_DIR}"
rm -rf "${GENERATED_DIR}"

echo "Generating definitions..."
go run ./cmd/generate "${GENERATED_DIR}"

echo "Compiling Go packages..."
go test ./...

echo "Validating generated CUE..."
for definition in "${GENERATED_DIR}"/component/*.cue; do
  "${VELA_BIN}" def vet "${definition}" >/dev/null
  echo "  valid: ${definition#"${ROOT_DIR}/"}"
done

echo "Dry-running valid applications..."
for application in \
  examples/queue-set.yaml \
  examples/managed-record.yaml \
  examples/file-share.yaml \
  examples/zone-replica.yaml \
  examples/composite-store.yaml \
  examples/tenant-space.yaml; do
  "${VELA_BIN}" dry-run --offline \
    -d "${GENERATED_DIR}/component" \
    -f "${application}" >/dev/null
  echo "  rendered: ${application}"
done

echo "Checking the expected validator failure..."
if validation_output=$("${VELA_BIN}" dry-run --offline \
  -d "${GENERATED_DIR}/component" \
  -f examples/zone-replica-invalid.yaml 2>&1); then
  echo "expected zone-replica-invalid.yaml to fail validation" >&2
  exit 1
fi

expected_message="zone 'us-east-1' must not match the primary zone"
if [[ "${validation_output}" != *"${expected_message}"* ]]; then
  echo "validator failed without the expected dynamic message" >&2
  printf '%s\n' "${validation_output}" >&2
  exit 1
fi

echo "  rejected with: ${expected_message}"
echo "Demo verification passed."
