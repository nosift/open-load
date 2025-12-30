# Open-Load (Custom)

This is a customized fork used for our deployment. Only differences from the original project are listed below.

Author: nosift
Repository: https://github.com/nosift/open-load

## Attribution

本项目基于 tbphp/gpt-load 二次开发并发布，原项目采用 MIT License。
原项目版权声明与许可文本见本仓库 LICENSE 文件。
本二开版本的新增/修改内容版权归本项目作者所有（如适用）

## Key Differences from Upstream

### UI/UX Enhancements
- **Apple-inspired UI**: Minimal visuals, refined typography, optimized spacing and shadows
- **Login page redesign**: Text-only logo treatment for cleaner appearance
- **Model Usage analytics**: New dedicated page with donut chart, hover tooltips, and flexible date-range filtering
- **Model usage stats**: Aggregated from request_logs with success/error/retry breakdown; auto-filters empty models
- **Key management UX**: Enhanced view-key modal, improved layout/scroll behavior, polished button/input styling
- **Manual validation**: Force-mark keys as valid/invalid; tabs display counts and default to showing valid keys
- **Status code tooltips**: Plain-language explanations in logs (e.g., 402 = insufficient balance/quota)
- **Sidebar optimization**: Left group list stays compact with internal scrolling; fixed right-side scroll behavior

### Backend Features (v1.1.0+)
- **Organization validation**: Auto-detect OpenAI organization verification from API headers
- **Premium model monitoring**: Track which keys can access high-tier models (e.g., gpt-5.2)
- **Smart key tracking**: Database fields for `is_organization_key`, `organization_id`, `organization_name`
- **Non-blocking monitoring**: Logs warnings for non-org keys without disrupting retry mechanism
- **Configurable premium models**: Set `premium_models` in system settings to monitor specific models
- **Enhanced visibility**: View organization status in admin UI and API responses

### Architecture Changes
- **Removed unified OpenAI aggregate**: Use per-channel/per-group proxy endpoints instead
- **Database migration v1.2.0**: Adds organization validation fields automatically on startup

For detailed organization validation features, see [ORGANIZATION_VALIDATION.md](ORGANIZATION_VALIDATION.md).

## Deployment (Docker)

```bash
cp .env.example .env
# Edit .env and set AUTH_KEY at minimum.

docker login ghcr.io
docker pull ghcr.io/nosift/open-load:latest
docker run -d --name open-load \
  -p 3001:3001 \
  --env-file .env \
  -v "$(pwd)/data:/app/data" \
  ghcr.io/nosift/open-load:latest
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
docker build -t open-load-custom .
docker run -d --name open-load \
  -p 3001:3001 \
  --env-file .env \
  -v "$(pwd)/data:/app/data" \
  open-load-custom
```

## Environment Variables (Common)

- AUTH_KEY: required, protects the UI and admin API.
- ENCRYPTION_KEY: optional, encrypts API keys at rest.
- DATABASE_DSN: optional, empty = SQLite at `./data/open-load.db`.
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
