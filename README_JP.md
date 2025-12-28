# GPT-Load（カスタム）

この二次開発版は自社運用向けで、元プロジェクトとの差分のみを記載しています。

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
# .env を編集し、AUTH_KEY、DATABASE_DSN、REDIS_DSN（必要なら）を設定。

docker build -t gpt-load-custom .
docker run -d --name gpt-load \
  -p 3001:3001 \
  --env-file .env \
  -v "$(pwd)/data:/app/data" \
  gpt-load-custom
```

## デプロイ（ソース）

```bash
cp .env.example .env
# .env を編集し、AUTH_KEY、DATABASE_DSN、REDIS_DSN（必要なら）を設定。

go mod tidy
make run
```
