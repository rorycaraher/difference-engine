#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

ACCOUNT_ID="${CLOUDFLARE_ACCOUNT_ID:?Set CLOUDFLARE_ACCOUNT_ID}"
IMAGE="registry.cloudflare.com/${ACCOUNT_ID}/difference-engine:latest"

echo "==> Building container image"
docker build -t "$IMAGE" "$REPO_ROOT/app"

echo "==> Pushing image to Cloudflare registry"
docker push "$IMAGE"

echo "==> Deploying Workers container"
cd "$REPO_ROOT/app" && wrangler deploy

echo "==> Deploying Pages site"
wrangler pages deploy "$REPO_ROOT/app/site" --project-name difference-engine

echo "Done."
