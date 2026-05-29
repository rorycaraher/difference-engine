#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

echo "==> Deploying Workers container"
cd "$REPO_ROOT/app" && npx wrangler deploy

echo "==> Deploying Pages site"
npx wrangler pages deploy "$REPO_ROOT/app/site" --project-name difference-engine

echo "Done."
