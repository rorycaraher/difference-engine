#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

echo "==> Building frontend"

rm -rf "$REPO_ROOT/frontend-dist"
cp -r "$REPO_ROOT/frontend" "$REPO_ROOT/frontend-dist"

ABOUT_TMP=$(mktemp)
trap 'rm -f "$ABOUT_TMP"' EXIT

pandoc -f markdown -t html "$REPO_ROOT/about.md" > "$ABOUT_TMP"

python3 - "$REPO_ROOT/frontend-dist/index.html" "$ABOUT_TMP" <<'PYEOF'
import sys

html_file = sys.argv[1]
about_file = sys.argv[2]

content = open(html_file).read()
about_html = open(about_file).read()
result = content.replace('<!-- ABOUT_CONTENT -->', about_html, 1)

with open(html_file, 'w') as f:
    f.write(result)
PYEOF

echo "==> Frontend built to frontend-dist/"
