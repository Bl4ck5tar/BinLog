#!/usr/bin/env bash
set -euo pipefail

API_BASE=${API_BASE:-http://localhost:8080/api}

if ! command -v k6 >/dev/null 2>&1; then
  echo "k6 not found. Install it first:"
  echo "  macOS (brew): brew install k6"
  exit 1
fi

echo "Running STRESS test against: $API_BASE"

OUT_DIR="$(dirname "$0")/results"
mkdir -p "$OUT_DIR"
TS=$(date +"%Y%m%d-%H%M%S")
SUMMARY="$OUT_DIR/summary-$TS.json"
RAW="$OUT_DIR/raw-$TS.json"

API_BASE="$API_BASE" k6 run "$(dirname "$0")/k6-stress.js" \
  --summary-export "$SUMMARY" \
  --out "json=$RAW"

echo "\nSummary JSON: $SUMMARY"

echo "Raw JSON: $RAW"
