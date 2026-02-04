#!/usr/bin/env bash
set -euo pipefail

API_BASE=${API_BASE:-http://localhost:8080/api}

if ! command -v k6 >/dev/null 2>&1; then
  echo "k6 not found. Install it first:"
  echo "  macOS (brew): brew install k6"
  echo "  or visit: https://k6.io/docs/get-started/installation/"
  exit 1
fi

echo "Running load test against: $API_BASE"
API_BASE="$API_BASE" k6 run "$(dirname "$0")/k6.js"
