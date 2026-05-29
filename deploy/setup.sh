#!/usr/bin/env bash
# Run once on the server as root to create the service user and directory layout.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

useradd --system --no-create-home --shell /usr/sbin/nologin difference-engine

install -d -m 755 /opt/difference-engine
install -d -m 755 /opt/difference-engine/frontend
install -d -m 755 -o difference-engine -g difference-engine /var/lib/difference-engine
install -d -m 755 -o difference-engine -g difference-engine /var/lib/difference-engine/output
install -d -m 700 -o root              -g root              /etc/difference-engine

cp "$SCRIPT_DIR/difference-engine.service" /etc/systemd/system/
systemctl daemon-reload
systemctl enable difference-engine

echo ""
echo "Setup complete."
echo "  1. Edit /etc/difference-engine/env  (use deploy/env.example as a template)"
echo "  2. Add the Caddy block from deploy/Caddyfile.snippet and reload Caddy"
echo "  3. Run scripts/deploy.sh to push the first build"
