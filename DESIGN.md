---
version: alpha
name: info-hub Dark
description: Personal RSS / X information dashboard. GitHub-dark-inspired theme built on Tailwind CSS v4 slate/sky/amber palette.
colors:
  background: "#0d1117"
  foreground: "#e6edf3"
  surface: "#0f172a"
  surface-muted: "#1e293b"
  border: "#1e293b"
  border-hover: "#475569"
  text-primary: "#f1f5f9"
  text-secondary: "#cbd5e1"
  text-muted: "#94a3b8"
  text-faint: "#64748b"
  primary: "#0284c7"
  accent: "#0284c7"
  on-accent: "#ffffff"
  favorite: "#fbbf24"
  warning: "#fbbf24"
  danger: "#f87171"
typography:
  h1:
    fontFamily: system-ui
    fontSize: 1.5rem
    fontWeight: 700
    letterSpacing: -0.025em
  body:
    fontFamily: system-ui
    fontSize: 0.875rem
    lineHeight: 1.5
  title:
    fontFamily: system-ui
    fontSize: 1rem
    fontWeight: 500
    lineHeight: 1.375
  caption:
    fontFamily: system-ui
    fontSize: 0.75rem
    lineHeight: 1.4
rounded:
  xs: 4px
  sm: 6px
  md: 8px
  full: 9999px
spacing:
  xs: 4px
  sm: 8px
  md: 12px
  lg: 16px
  xl: 24px
components:
  card:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.text-primary}"
    rounded: "{rounded.md}"
    padding: 12px
  chip-active:
    backgroundColor: "{colors.accent}"
    textColor: "{colors.on-accent}"
    rounded: "{rounded.full}"
    padding: 6px
  chip-inactive:
    backgroundColor: "{colors.surface-muted}"
    textColor: "{colors.text-secondary}"
    rounded: "{rounded.full}"
    padding: 6px
  input:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.foreground}"
    rounded: "{rounded.sm}"
    padding: 8px
  button-ghost:
    backgroundColor: "{colors.background}"
    textColor: "{colors.text-muted}"
    rounded: "{rounded.sm}"
    padding: 8px
  badge:
    backgroundColor: "{colors.surface-muted}"
    textColor: "{colors.text-muted}"
    rounded: "{rounded.xs}"
    padding: 4px
  button-fav:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.favorite}"
    rounded: "{rounded.xs}"
    padding: 4px
  error-banner:
    backgroundColor: "#450a0a"
    textColor: "{colors.danger}"
    rounded: "{rounded.xs}"
    padding: 8px
  warning-banner:
    backgroundColor: "#451a03"
    textColor: "{colors.warning}"
    rounded: "{rounded.xs}"
    padding: 6px
---

# info-hub Design System

## Overview

info-hub は個人用 RSS / X 情報ダッシュボード。GitHub ダークテーマに寄せた配色で、
長時間の流し読みでも目が疲れない低コントラストのダーク UI を基本とする。
実装は Tailwind CSS v4 のユーティリティクラス（slate / sky / amber パレット）を用いる。
ライトテーマは提供しない（`color-scheme: dark` 固定）。

## Colors

- `background` (#0d1117): ページ全体の背景。GitHub dark と同一
- `foreground` (#e6edf3): 基本テキスト色
- `surface` (#0f172a, slate-900): カード・入力欄の背景。既読カードは opacity 50% で沈める
- `surface-muted` (#1e293b, slate-800): チップ・バッジ・ホバー背景
- `border` (#1e293b) / `border-hover` (#475569): カード境界。ホバーで明るくして操作可能性を示す
- テキスト階層: `text-primary` (slate-100) > `text-secondary` (slate-300) > `text-muted` (slate-400) > `text-faint` (slate-500)
- `accent` (#0284c7, sky-600): 選択状態・フォーカスリング・アクションの単一アクセント
- `favorite` (#fbbf24, amber-400): お気に入り（★）専用。警告メッセージにも流用
- `danger` (#f87171, red-400): データ読み込みエラー等の失敗表示

## Typography

- フォントは `system-ui`（Tailwind デフォルトのシステムフォントスタック）。Web フォントは読み込まない
- `h1`: ページタイトルのみ（text-2xl / bold / tracking-tight）
- `title`: 記事タイトル（font-medium / leading-snug）。未読は text-primary、既読は text-muted に落とす
- `body`: UI ラベル・ボタン（text-sm）
- `caption`: メタ情報・バッジ・補足（text-xs）
- 本文は日本語主体のため字間調整はせず、h1 のみ tracking-tight

## Layout

- コンテンツは `max-w-3xl` の単一カラム、中央寄せ（`mx-auto px-4 py-6`）
- リストは `space-y-2`、フィルタ行は `flex flex-wrap gap-2〜3` で折り返し対応
- スペーシングは 4px 基調のスケール（xs=4 / sm=8 / md=12 / lg=16 / xl=24）
- カード内は `px-4 py-3`、チップ・入力欄は `px-3 py-1.5` を標準とする

## Elevation & Depth

- box-shadow は使わない。階層は背景色の明度差（background < surface < surface-muted）と
  境界線（border → border-hover）で表現する
- 既読アイテムは `opacity-50` で背面に沈める
- オーバーレイ・モーダルは現状なし。追加する場合も影ではなく背景の明度差で表現する

## Shapes

- カード・エラー表示: `rounded.md` (8px, rounded-lg)
- 入力欄・ボタン: `rounded.sm` (6px, rounded-md)
- カテゴリバッジ: `rounded.xs` (4px, rounded)
- フィルタチップ: `rounded.full`（ピル型）
- シャープな角（角丸なし）は使わない

## Components

- **card** (`ItemCard`): surface 背景 + border。ホバーで border-hover。既読で opacity-50
- **chip-active / chip-inactive**: カテゴリフィルタ。選択中のみ accent 背景、それ以外は surface-muted
- **input**: 検索欄。フォーカスで border を accent に変更（`focus:border-sky-600 focus:outline-none`）
- **button-ghost**: 「表示中を既読に」等の低優先アクション。枠線のみでホバー時に surface-muted
- **badge**: カテゴリ表示。surface-muted 背景の小さな矩形
- チェックボックスは `accent-sky-600`（★フィルタのみ `accent-amber-500`）

## Do's and Don'ts

- Do: 色は必ず上記トークン（Tailwind の slate / sky / amber クラス）から選ぶ
- Do: 状態変化は `transition-colors` を付ける
- Do: 新しい強調色が必要なときは accent（sky）を使い回す
- Don't: `#fff` や任意の hex をクラスに直書きしない（`bg-[#xxx]` 禁止）
- Don't: box-shadow・グラデーション・ライトテーマ用の色を導入しない
- Don't: sky / amber / red 以外の有彩色（green, purple 等)を新規追加しない — 追加時はこのファイルを先に更新する
