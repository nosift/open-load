# GPT-Load (Custom)

This is a customized fork used for our deployment. Only differences from the original project are listed below.

## Differences

- UI overhaul (Apple-like minimal visuals, typography, spacing, and shadows).
- Login page refresh with a text-only logo treatment.
- New Model Usage page with donut chart, hover tooltip, and free date-range filter.
- Model usage stats aggregated from request_logs with success/error/retry counts; empty models filtered.
- Key management UX polish (view-key modal, layout/scroll fixes, button/input styling).
- Removed unified OpenAI aggregate proxy; use per-channel/per-group proxy endpoints.

## Deployment (Docker)

```bash
cp .env.example .env
# Edit .env and set AUTH_KEY at minimum.

docker login ghcr.io
docker pull ghcr.io/nosift/gpt-load:latest
docker run -d --name gpt-load \
  -p 3001:3001 \
  --env-file .env \
  -v "$(pwd)/data:/app/data" \
  ghcr.io/nosift/gpt-load:latest
```

## Image Tags & Update Policy

- `latest`: built automatically on every push to `main` (continuous updates).
- `vX.Y.Z-custom`: built when a git tag like `v1.0.1-custom` is pushed (stable versions).

Update a running service:

```bash
docker compose pull
docker compose up -d
```

## Release a Stable Version (Tag)

```bash
# After committing changes on main:
git tag v1.0.1-custom
git push origin v1.0.1-custom
```

### Build Locally

```bash
docker build -t gpt-load-custom .
docker run -d --name gpt-load \
  -p 3001:3001 \
  --env-file .env \
  -v "$(pwd)/data:/app/data" \
  gpt-load-custom
```

## Environment Variables (Common)

- AUTH_KEY: required, protects the UI and admin API.
- ENCRYPTION_KEY: optional, encrypts API keys at rest.
- DATABASE_DSN: optional, empty = SQLite at `./data/gpt-load.db`.
- REDIS_DSN: optional, empty = in-memory cache.
- TZ: timezone, default `Asia/Shanghai`.
- PORT / HOST: server bind.
- ENABLE_CORS / ALLOWED_ORIGINS / ALLOWED_METHODS / ALLOWED_HEADERS / ALLOW_CREDENTIALS.
- LOG_LEVEL / LOG_FORMAT / LOG_ENABLE_FILE / LOG_FILE_PATH.
- MAX_CONCURRENT_REQUESTS.

Full list: `.env.example`.

## Deployment (Source, no Docker)

```bash
cp .env.example .env
# Edit .env and set AUTH_KEY at minimum.

# Prereqs: Go 1.24+ toolchain and Node.js 18+.
cd web && npm install && npm run build
cd ..
go mod tidy
go run ./main.go
```
