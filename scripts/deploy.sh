#!/usr/bin/env bash
# Deploy to the Hetzner server.
#
# Required env vars (set in .env or export before running):
#   DEPLOY_HOST   e.g. app@your-server.example.com
#
# Optional:
#   GOARCH        target CPU arch (default: amd64; use arm64 for Hetzner CAX)
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

if [ -f "$REPO_ROOT/.env" ]; then
  set -a && source "$REPO_ROOT/.env" && set +a
fi

: "${DEPLOY_HOST:?DEPLOY_HOST is not set}"
GOARCH="${GOARCH:-amd64}"
BINARY="$REPO_ROOT/difference-engine"

echo "==> Building linux/$GOARCH binary"
cd "$REPO_ROOT/backend"
GOOS=linux GOARCH="$GOARCH" CGO_ENABLED=0 go build -o "$BINARY" .

echo "==> Uploading binary"
rsync -az "$BINARY" "$DEPLOY_HOST:/opt/difference-engine/difference-engine"

echo "==> Uploading frontend"
rsync -az --delete "$REPO_ROOT/frontend/" "$DEPLOY_HOST:/opt/difference-engine/frontend/"

echo "==> Restarting service"
ssh "$DEPLOY_HOST" 'sudo systemctl restart difference-engine && sudo systemctl status difference-engine --no-pager'

echo "==> Done"
rm "$BINARY"
