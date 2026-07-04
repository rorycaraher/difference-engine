#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BACKEND_DIR="$REPO_ROOT/backend"

if [ -f "$REPO_ROOT/.env" ]; then
  set -a && source "$REPO_ROOT/.env" && set +a
fi

"$REPO_ROOT/scripts/build-frontend.sh"

cd "$BACKEND_DIR"
SITE_DIR="$REPO_ROOT/frontend-dist" go run .
