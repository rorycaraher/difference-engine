# difference-engine

Generative music player for "NLTL – First Principles". Each visit randomly selects 3 stems from 24 available, plays them in the browser via Web Audio API, and lets the user download a server-side ffmpeg mixdown as MP3. Live at `de.nothinglefttolearn.com`.

## Repository layout

| Path | Purpose |
|---|---|
| `backend/` | Go HTTP server — routes, config, request handlers |
| `backend/mixer/` | ffmpeg-based audio mixing; abstracts local filesystem vs Cloudflare R2 |
| `backend/store/` | SQLite persistence for mixdown requests (stems + volumes) |
| `frontend/` | Static HTML/CSS/JS — Web Audio API playback + download UI |
| `deploy/` | Systemd unit file, Caddyfile snippet, server setup script |
| `scripts/` | `dev.sh` (local dev server), `deploy.sh` (manual deploy to VPS) |
| `.github/workflows/` | CI/CD — build + rsync + service restart on push to `main` |

## Local development

**Prerequisites:** Go 1.25+, ffmpeg, a directory of stem MP3 files (or R2 credentials).

```bash
cp .env.example .env          # fill in STEMS_DIR and OUTPUT_DIR at minimum
./scripts/dev.sh              # starts Go server on :8000, serves frontend/
```

No frontend build step — plain JS, served directly from `frontend/`.

Run backend tests:
```bash
cd backend && go test ./...
```

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8000` | HTTP listen port |
| `STEMS_DIR` | — | Base directory for local stem files |
| `OUTPUT_DIR` | — | Directory where ffmpeg writes mixdown MP3s |
| `DB_PATH` | `de.db` | SQLite database file path |
| `SITE_DIR` | — | Serve frontend from this directory (dev only) |
| `R2_ACCOUNT_ID` | — | Cloudflare account ID |
| `R2_ACCESS_KEY_ID` | — | R2 access key |
| `R2_SECRET_ACCESS_KEY` | — | R2 secret key |
| `R2_STEMS_BUCKET` | — | R2 bucket name |

**Storage toggle:** if `R2_STEMS_BUCKET` is set, the entire app uses R2 for stems (list, fetch, presign). If unset, it uses the local filesystem under `STEMS_DIR`. There is no mixed mode.

See `.env.example` for local dev and `deploy/env.example` for the production systemd EnvironmentFile.

## API endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/stems/{track}/count` | List available stems for a track |
| `GET` | `/stems/{track}/{stem}` | Serve a stem — redirects to presigned URL (R2) or serves file (local) |
| `POST` | `/mixdown` | Mix stems at given volumes; returns MP3 + `X-Mixdown-ID` header |
| `GET` | `/mixdown/{id}` | Recall and re-mix a previously recorded request |

**POST /mixdown body:**
```json
{
  "track": "track-name",
  "de_values": {
    "stems": ["stem1.mp3", "stem2.mp3"],
    "volumes": [0.8, 0.6]
  }
}
```

## Architecture notes

**R2 vs local:** `mixer.go` checks `m.cfg.R2StemsBucket != ""` on every list/fetch call. Add the four `R2_*` vars to switch modes; remove them to go back to local.

**ffmpeg mixing:** `buildFFmpegArgs()` in `mixer/mixer.go` constructs a `filter_complex` chain — one `volume` filter per input stream, then `amix`, then `loudnorm`. ffmpeg runs as a subprocess via `exec.Command`; stdout/stderr pipe to the process's own stdout/stderr for live logging.

**SQLite:** one open connection (`db.SetMaxOpenConns(1)`) + WAL mode to avoid writer contention. Stems and volumes are stored as JSON strings (not typed array columns). Schema migrations are version-gated via `PRAGMA user_version`.

**Path safety:** `safeSegment` regex in `main.go` validates `{track}` and `{stem}` path values before use to prevent directory traversal.

**Handler convention:** all HTTP handlers are methods on `*server` named `handle<X>()`. Error responses use `log.Printf` + `http.Error`. `fmt.Errorf("context: %w", err)` is used throughout for error wrapping.

## Deployment

**Automated (preferred):** push to `main` → GitHub Actions builds a static Linux amd64 binary, rsyncs it and `frontend/` to the Hetzner VPS, restarts the systemd service.

**Manual:**
```bash
./scripts/deploy.sh                        # amd64 (default)
GOARCH=arm64 ./scripts/deploy.sh           # arm64
```

**Server stack:** Caddy handles TLS and reverse-proxies `/stems*` and `/mixdown*` to the Go server on `:8000`; all other paths are served as static files from `/opt/difference-engine/frontend/`. See `deploy/Caddyfile.snippet`.

**Production env file:** `/etc/difference-engine/env` (mode 600), sourced by the systemd unit.

## Common commands

```bash
# Local dev server
./scripts/dev.sh

# Backend tests
cd backend && go test ./...

# Build static Linux binary manually
cd backend && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../difference-engine .

# Deploy manually
./scripts/deploy.sh
```
