# GPT-Load（二开）

这是用于我方部署的二开版本，仅保留与原项目的差异。

## 差异

- UI 统一视觉重构（更简约的排版、间距、阴影与字体）。
- 登录页与文字 Logo 调整。
- 新的模型统计页面（环形占比、悬浮提示、自由日期范围）。
- 模型统计基于 request_logs 聚合，补充成功/错误/重试，过滤空模型。
- 密钥管理体验优化（查看密钥弹窗、布局/滚动修复、按钮/输入样式优化）。
- 移除统一 OpenAI 聚合接口，恢复按渠道/分组访问。

## 部署（Docker）

```bash
cp .env.example .env
# 编辑 .env，至少设置 AUTH_KEY。

docker login ghcr.io
docker pull ghcr.io/nosift/gpt-load:latest
docker run -d --name gpt-load \
  -p 3001:3001 \
  --env-file .env \
  -v "$(pwd)/data:/app/data" \
  ghcr.io/nosift/gpt-load:latest
```

## 镜像标签与更新策略

- `latest`：每次 `main` 分支更新都会自动构建（持续更新）。
- `vX.Y.Z-custom`：推送对应 tag（如 `v1.0.1-custom`）时构建（稳定版本）。

更新已运行服务：

```bash
docker compose pull
docker compose up -d
```

### 本地构建

```bash
docker build -t gpt-load-custom .
docker run -d --name gpt-load \
  -p 3001:3001 \
  --env-file .env \
  -v "$(pwd)/data:/app/data" \
  gpt-load-custom
```

## 环境变量（常用）

- AUTH_KEY：必填，用于保护管理端与 API。
- ENCRYPTION_KEY：可选，用于加密保存密钥。
- DATABASE_DSN：可选，留空则使用 `./data/gpt-load.db`（SQLite）。
- REDIS_DSN：可选，留空则使用内存缓存。
- TZ：时区，默认 `Asia/Shanghai`。
- PORT / HOST：服务绑定地址。
- ENABLE_CORS / ALLOWED_ORIGINS / ALLOWED_METHODS / ALLOWED_HEADERS / ALLOW_CREDENTIALS。
- LOG_LEVEL / LOG_FORMAT / LOG_ENABLE_FILE / LOG_FILE_PATH。
- MAX_CONCURRENT_REQUESTS。

完整列表见：`.env.example`。

## 部署（源码，无 Docker）

```bash
cp .env.example .env
# 编辑 .env，至少设置 AUTH_KEY。

# 依赖：Go 1.24+ toolchain 与 Node.js 18+。
cd web && npm install && npm run build
cd ..
go mod tidy
go run ./main.go
```
