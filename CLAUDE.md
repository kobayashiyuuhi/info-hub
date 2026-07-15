# info-hub

個人用 RSS / X 情報ダッシュボード。Next.js 15 (App Router, static export) + Tailwind CSS v4 + Go fetcher。

## デザインルール

- **UI 生成・変更時は必ずルートの `DESIGN.md` のデザイントークンを参照すること**
- 色は DESIGN.md の colors トークン（Tailwind の slate / sky / amber クラス）から選ぶ。hex 直書き（`bg-[#xxx]` 等）禁止
- 新しい色・形状・コンポーネントスタイルを導入する場合は、先に DESIGN.md を更新してから実装する
- DESIGN.md 変更後は lint で検証する:
  `node node_modules/@google/design.md/dist/index.js lint DESIGN.md`
  （エラー 0 を維持。alpha 仕様のためスキーマは変動し得る）

## 構成

- `app/` / `components/` — Next.js フロントエンド（ダークテーマ固定）
- `fetcher/` — Go 製フィード取得バッチ（`public/data/*.json` を生成）
- `feeds.yml` — 購読フィード定義
