#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

ENV_FILE=""
if [ -f "$REPO_ROOT/.env" ]; then
  ENV_FILE="--env-file $REPO_ROOT/.env"
fi

docker run --rm -it \
  $ENV_FILE \
  -p 8000:8000 \
  difference-engine:latest
