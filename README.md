# info-hub

自分専用の情報収集ダッシュボード（Feedly代替）。RSS/Atomフィードと X（ブリッジ経由）の最新記事を一括表示する。GitHub Pages でホスト。

## 構成

```
┌─────────────┐   cron (6時間ごと)
│ GitHub       │  ┌────────────────────────────────┐
│ Actions      │  │ fetch.yml                       │
│              │  │  Go fetcher が feeds.yml を読み  │
│              │  │  RSS/Atom を並行取得             │
│              │  │  → public/data/*.json を生成     │
│              │  │  → 差分あれば commit & push      │
│              │  │  → deploy.yml を起動             │
│              │  └────────────────────────────────┘
│              │  ┌────────────────────────────────┐
│              │  │ deploy.yml (main push / dispatch)│
│              │  │  Next.js 静的エクスポート        │
│              │  │  → GitHub Pages にデプロイ       │
│              │  └────────────────────────────────┘
└─────────────┘
       ↓
  https://kobayashiyuuhi.github.io/info-hub/
  （フロントは静的JSONを読むだけ。既読・★・検索は全部 localStorage / クライアントサイド）
```

- **fetcher/** — Go 製フェッチャー（gofeed）。タイムアウト付き並行取得、失敗フィードはスキップして `meta.json` にエラー記録。全アイテムをマージ・日付降順・直近500件を `public/data/items.json` に出力
- **app/ + components/** — Next.js (App Router, TypeScript, `output: 'export'`, Tailwind CSS)。ダークテーマ。カテゴリタブ / キーワード検索 / 既読グレーアウト / 未読のみ表示 / ★お気に入り
- **feeds.yml** — 購読リスト（このファイルを編集するだけでフィード追加・削除）

## フィードの追加方法

`feeds.yml` にエントリを追加して main に push するだけ。次回の fetch 実行（6時間ごと、または手動）で反映される。

```yaml
feeds:
  - name: 表示名
    url: https://example.com/rss
    category: tech   # カテゴリタブになる。自由に追加可（tech / ai / x など）
```

すぐ反映したい場合は Actions タブから **Fetch feeds** を `Run workflow` で手動実行する。

## X（Twitter）の購読

X には公式RSSがないため、RSSブリッジのURLを普通のRSSとして `feeds.yml` に登録する。

```yaml
  - name: "@OpenAI (X)"
    url: https://rsshub.app/twitter/user/OpenAI
    category: x
```

### 注意（ブリッジは不安定）

- **RSSHub 公式インスタンス（rsshub.app）はレート制限・ブロックで頻繁に失敗する。** セルフホスト（Docker / Vercel 等）か、動作している公開ミラーを使うのが現実的
- **Nitter** 系インスタンス（`https://<nitter-host>/<user>/rss`）も候補だが、稼働インスタンスが頻繁に入れ替わる
- 失敗したフィードはスキップされ、ダッシュボード上部に「取得失敗」として表示される（他のフィードには影響しない）

## ローカル開発

前提: Go 1.24+ / Node.js 22+

```bash
# 1. データ生成（fetcher ディレクトリから実行）
cd fetcher
go run . -feeds ../feeds.yml -out ../public/data
go test ./...   # ユニットテスト

# 2. フロント
cd ..
npm install
npm run dev     # http://localhost:3000 （開発時は basePath なし）
npm run build   # 静的エクスポート → out/
```

本番ビルドは `basePath: '/info-hub'` が付く（GitHub Pages のサブパス対応、`next.config.mjs` 参照）。

## 出力データ形式

- `public/data/items.json` — 記事配列 `{id, title, url, source, category, published, summary}`
- `public/data/meta.json` — `{generatedAt, feeds: [{name, url, category, status, error?, itemCount}]}`

## セットアップ（初回のみ・手動手順）

GitHub Pages が未設定の場合:

1. リポジトリの **Settings → Pages → Build and deployment → Source** を **GitHub Actions** にする

   CLI の場合:
   ```bash
   gh api -X POST repos/kobayashiyuuhi/info-hub/pages -f build_type=workflow
   ```
2. Actions タブで **Deploy to GitHub Pages** を手動実行（または main に push）
3. `https://kobayashiyuuhi.github.io/info-hub/` で確認
