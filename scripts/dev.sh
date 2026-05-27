#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
APP_DIR="$REPO_ROOT/app"

if [ -f "$REPO_ROOT/.env" ]; then
  set -a && source "$REPO_ROOT/.env" && set +a
fi

cd "$APP_DIR"
uv run fastapi dev main.py --port 8000
