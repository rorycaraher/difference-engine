#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
IMAGE="${IMAGE:-difference-engine:latest}"

echo "==> Building image"
docker build \
  --file "$REPO_ROOT/backend/Dockerfile" \
  --tag "$IMAGE" \
  "$REPO_ROOT"

echo "==> Pushing image"
docker push "$IMAGE"

echo "Done."
