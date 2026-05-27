#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

docker build \
  --tag difference-engine:latest \
  "$REPO_ROOT/app"
