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
# Edit .env and set AUTH_KEY, DATABASE_DSN, REDIS_DSN if needed.

docker login ghcr.io
docker pull ghcr.io/nosift/gpt-load:v1.0.0-custom
docker run -d --name gpt-load \
  -p 3001:3001 \
  --env-file .env \
  -v "$(pwd)/data:/app/data" \
  ghcr.io/nosift/gpt-load:v1.0.0-custom
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

## Deployment (Source)

```bash
cp .env.example .env
# Edit .env and set AUTH_KEY, DATABASE_DSN, REDIS_DSN if needed.

go mod tidy
make run
```
