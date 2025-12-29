# Open-Load（二开）

这是用于我方部署的二开版本，仅保留与原项目的差异。

作者：nosift  
仓库主页：https://github.com/nosift/gpt-load

## 版权与许可

本项目基于 tbphp/gpt-load 二次开发并发布，原项目采用 MIT License。
原项目版权声明与许可文本见本仓库 LICENSE 文件。
本二开版本的新增/修改内容版权归本项目作者所有（如适用）

## 差异

### UI/交互

- UI 统一视觉重构（更简约的排版、间距、阴影与字体）。
- 顶部/底部/登录页的文字 Logo 与信息重构（本仓库 UI 名称为 Open Load）。
- 分组页布局优化：左侧分组栏固定不动；左栏保持“短”，列表超出时内部滚动；右侧主内容区域滚动。
- 分组列表图标使用稳定缩写（如 AI/GE/AN/AG），避免在容器/字体缺失场景出现“图标消失/乱码”。

### 功能

- 新的模型统计页面（环形占比、悬浮提示、自由日期范围）。
- 模型统计基于 request_logs 聚合，补充成功/错误/重试，过滤空模型。
- 密钥管理体验优化（查看密钥弹窗、布局/滚动修复、按钮/输入样式优化）。
- 手动验证后强制标记有效/无效，并在状态标签页展示数量；默认展示“有效”密钥。
- 密钥工具栏将“更多”改为“验证”下拉入口（批量验证/导出/清理等）。
- 日志页面：状态码支持悬浮/详情“大白话解释”（如 402=余额/额度不足/需要付费或配额用完）。
- 移除统一 OpenAI 聚合接口，恢复按渠道/分组访问。

### 更新与版本

- 前端版本检查：优先读取 GitHub Release；无 Release 时自动回退到 tags 获取最新版本（避免“检查失败”）。
- GitHub Actions：
  - push 到 `main` 自动构建并推送 `:latest` 镜像（持续更新）。
  - push tag（如 `v1.0.1-custom`）构建并推送同名 tag 镜像（稳定版本）。
  - tag 也会触发多平台二进制打包的草稿 Release（如需对外发布可在 GitHub 上手动发布草稿）。

## 常见状态码（大白话）

| 状态码 | 含义（大白话） |
| --- | --- |
| 401 | Key 无效/过期，或鉴权没带对。 |
| 402 | 余额/额度不足（需要付费或配额用完，计费渠道常见）。 |
| 403 | 被拒绝/无权限（Key 权限不足、风控、来源限制等）。 |
| 429 | 被限流/请求太频繁（降低并发或稍后重试）。 |
| 5xx | 上游/网关异常（多为服务端问题，稍后重试）。 |

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

## 发布稳定版本（打 Tag）

发布前检查清单（建议每次都走一遍）：

- 确认功能/UI 已自测通过（尤其是密钥页、日志页、容器部署场景）。
- 同步更新主页文档：`README_CN.md`（必要时也更新 `README.md`）。
- `git status` 确认没有未提交修改。
- 合并/提交后先推送到 `main`，确认 GitHub Actions 的 `latest` 镜像构建成功。
- 再打 tag 发布稳定版（会触发 tag 镜像构建；版本检查也会更容易识别更新）。

```bash
# main 上提交改动后：
git tag v1.0.1-custom
git push origin v1.0.1-custom
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
