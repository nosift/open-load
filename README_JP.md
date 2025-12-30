# Open-Load（カスタム）

この二次開発版は自社運用向けで、元プロジェクトとの差分のみを記載しています。

作者: nosift  
リポジトリ: https://github.com/nosift/open-load

## ライセンス

本项目基于 tbphp/gpt-load 二次开发并发布，原项目采用 MIT License。
原项目版权声明与许可文本见本仓库 LICENSE 文件。
本二开版本的新增/修改内容版权归本项目作者所有（如适用）

## 差分

- UI を全体的に再設計（よりミニマルなレイアウト、余白、影、タイポグラフィ）。
- ログイン画面とテキストロゴの刷新。
- 新しいモデル統計ページ（ドーナツチャート、ホバー表示、自由な日付範囲）。
- request_logs から集計し、成功/エラー/リトライを追加、空のモデル名を除外。
- キー管理 UX 改善（キー閲覧モーダル、レイアウト/スクロール修正、ボタン/入力の調整）。
- 統一 OpenAI 集約インターフェースを削除し、チャンネル/グループ別のアクセスに戻しました。

## デプロイ（Docker）

```bash
cp .env.example .env
# .env を編集し、最低でも AUTH_KEY を設定。

docker login ghcr.io
docker pull ghcr.io/nosift/open-load:latest
docker run -d --name open-load \
  -p 3001:3001 \
  --env-file .env \
  -v "$(pwd)/data:/app/data" \
  ghcr.io/nosift/open-load:latest
```

## イメージタグと更新方針

- `latest`：`main` への push ごとに自動ビルド（継続更新）。
- `vX.Y.Z-custom`：`v1.0.1-custom` のような tag push でビルド（安定版）。

稼働中サービスの更新：

```bash
docker compose pull
docker compose up -d
```

## 安定版のリリース（Tag）

```bash
# main にコミット後：
git tag v1.0.1-custom
git push origin v1.0.1-custom
```

### ローカルビルド

```bash
docker build -t open-load-custom .
docker run -d --name open-load \
  -p 3001:3001 \
  --env-file .env \
  -v "$(pwd)/data:/app/data" \
  open-load-custom
```

## 環境変数（主な項目）

- AUTH_KEY：必須、管理 UI と API を保護。
- ENCRYPTION_KEY：任意、キーを暗号化して保存。
- DATABASE_DSN：任意、空なら `./data/open-load.db`（SQLite）。
- REDIS_DSN：任意、空ならメモリキャッシュ。
- TZ：タイムゾーン、既定 `Asia/Shanghai`。
- PORT / HOST：サーバーバインド。
- ENABLE_CORS / ALLOWED_ORIGINS / ALLOWED_METHODS / ALLOWED_HEADERS / ALLOW_CREDENTIALS。
- LOG_LEVEL / LOG_FORMAT / LOG_ENABLE_FILE / LOG_FILE_PATH。
- MAX_CONCURRENT_REQUESTS。

全項目は `.env.example` を参照。

## デプロイ（ソース、Docker なし）

```bash
cp .env.example .env
# .env を編集し、最低でも AUTH_KEY を設定。

# 依存関係：Go 1.24+ toolchain と Node.js 18+。
cd web && npm install && npm run build
cd ..
go mod tidy
go run ./main.go
```
