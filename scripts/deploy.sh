#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
IMAGE="${IMAGE:-difference-engine:latest}"

echo "==> Building image"
docker build --tag "$IMAGE" "$REPO_ROOT/app"

echo "==> Pushing image"
docker push "$IMAGE"

echo "Done."
