# Difference Engine

A generative music player for [NLTL – First Principles](https://de.nothinglefttolearn.com).

Each page load selects 3 stems at random from a pool of 24, plays them simultaneously at randomised volumes, and lets you download the unique mix as an MP3.

## How it works

- **Frontend** — plain HTML/CSS/JS. On load it fetches the list of available stems from the backend, shuffles them, and streams three into Web Audio simultaneously.
- **Backend** — a single Go binary that serves the API and (in local dev) the static frontend.
  - `GET /stems/{track}/count` — returns the list of stem identifiers for a track, read from the local filesystem or Cloudflare R2.
  - `GET /stems/{track}/{stem}` — streams a stem directly (local) or redirects to a short-lived R2 presigned URL.
  - `POST /mixdown` — mixes the chosen stems at the given volumes via ffmpeg and returns the resulting MP3.
  - `GET /mixdown/{id}` — replays a previously recorded mix.
- **Stems** — stored in Cloudflare R2 in production, or a local directory in development.
- **Database** — SQLite records every mixdown request so mixes can be recalled by ID.

## Local development

**Prerequisites:** Go, ffmpeg, and either a local stems directory or R2 credentials.

```
cp .env.example .env
# edit .env — set STEMS_DIR to a folder containing first-principles/*.mp3
scripts/dev.sh
```

The server starts on port 8000 and serves the frontend from `frontend/`.

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8000` | HTTP listen port |
| `SITE_DIR` | `frontend` | Static files to serve (leave unset in production — Caddy handles it) |
| `DB_PATH` | `de.db` | SQLite database path |
| `STEMS_DIR` | — | Base directory for local stems (`STEMS_DIR/track/stem.mp3`) |
| `OUTPUT_DIR` | `./output` | Where ffmpeg writes mixdown files |
| `R2_ACCOUNT_ID` | — | Cloudflare account ID |
| `R2_ACCESS_KEY_ID` | — | R2 API token key ID |
| `R2_SECRET_ACCESS_KEY` | — | R2 API token secret |
| `R2_STEMS_BUCKET` | — | R2 bucket name — presence of this switches stems from local to R2 |

## Deployment

Targets a Hetzner VPS running Caddy + systemd. Set `DEPLOY_HOST` in `.env` (e.g. `app@your-server.example.com`), then:

```
scripts/deploy.sh
```

This cross-compiles a static Linux binary, rsyncs it and the frontend to `/opt/difference-engine/`, and restarts the service. Set `GOARCH=arm64` for Hetzner CAX (Ampere) instances.

Caddy proxies `/stems*` and `/mixdown*` to the Go backend and serves everything else as static files from `/opt/difference-engine/frontend/`. A sample Caddyfile snippet and systemd unit are in `deploy/`.

## Adding stems

Drop `.mp3` files into the track subdirectory (locally: `STEMS_DIR/first-principles/`, in R2: `first-principles/`). The backend lists available files at startup of each request — no config change needed. Avoid uploading macOS `._*` metadata files; the backend filters them out but they'll clutter the bucket.
