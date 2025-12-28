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
# 编辑 .env，设置 AUTH_KEY、DATABASE_DSN、REDIS_DSN（如需）。

docker build -t gpt-load-custom .
docker run -d --name gpt-load \
  -p 3001:3001 \
  --env-file .env \
  -v "$(pwd)/data:/app/data" \
  gpt-load-custom
```

## 部署（源码）

```bash
cp .env.example .env
# 编辑 .env，设置 AUTH_KEY、DATABASE_DSN、REDIS_DSN（如需）。

go mod tidy
make run
```
