# design-template1

Go 標準ライブラリだけで動く管理画面向けの **デザインテンプレート** です。フロントは **Tailwind CSS v4** と **Preline UI** でコンポーネントの見た目と挙動をそろえています。

## 技術スタック

| 領域 | 内容 |
|------|------|
| バックエンド | Go 1.21+、`net/http`、`html/template`（外部 Go 依存なし） |
| スタイル | Tailwind CSS v4、`@tailwindcss/cli` |
| コンポーネント基盤 | [Preline UI](https://preline.co/)（`preline` npm パッケージ） |
| フォーム用 | `@tailwindcss/forms`（Preline 推奨） |

インタラクティブな部品（ドロップダウン、モーダル、ツールチップなど）は Preline の `preline.js`（ビルド時に `static/js/` へコピー）で初期化されます。

## 画面仕様（デモページ）

`GET /` で **1 枚の管理画面デモ**を表示します。

- **レイアウト**: 左サイドバー（ナビ）＋上部バー（タイトル・ユーザードロップダウン）＋メインエリア
- **スタック表記**: Go / Tailwind CSS / Preline UI（**ソリッド**のバッジ表現）
- **アラート**: 情報・成功・注意・エラー・ヒント（見た目のバリエーション）
- **ツールチップ**: Preline `HSTooltip`（`[--placement:*]`、`[--trigger:hover|click|focus]`、`[--scope:window]`）。**矢印**は Preline 非搭載のため、`static/css/input.css` で `data-placement` に応じた `::after` を定義
- **ボタン**: ベーシック／カラー（ソリッド）／ソフト／アイコン付き／アイコンのみ／サイズ／ゴースト・リンク
- **テーブル**: ユーザー一覧サンプル（`main.go` から渡すダミーデータ）、ステータスはソリッドバッジ
- **モーダル**: Preline `HSOverlay`（`data-hs-overlay`）

テンプレートは `templates` 以下の `.html` を再帰的に読み込み、`{{define "admin"}}` を `ExecuteTemplate(..., "admin", data)` で描画します。

## ディレクトリ構成

```
design-template1/
├── main.go                 # HTTP サーバ・ルーティング
├── go.mod
├── package.json            # npm スクリプト・Tailwind / Preline 依存
├── scripts/
│   └── copy-preline.js     # preline.js を static/js へコピー
├── static/
│   ├── css/
│   │   ├── input.css       # Tailwind エントリ（@source / Preline / ツールチップ矢印など）
│   │   └── output.css      # ビルド成果物（リポジトリに含める想定）
│   └── js/
│       └── preline.js      # npm run build 時に生成（コピー）
├── templates/
│   ├── admin.html
│   └── partials/
│       ├── sidebar.html
│       └── topbar.html
├── .vscode/
│   ├── launch.json         # 実行とデバッグ用
│   └── tasks.json          # npm: build タスク
└── README.md
```

## 必要な環境

- Go 1.21 以上
- Node.js / npm（CSS ビルドと Preline JS のコピー用）

## セットアップと起動

```bash
cd design-template1
npm install
npm run build    # Tailwind → output.css、preline.js を static/js にコピー
go run .
```

ブラウザで `http://localhost:8080/` を開きます。

### ポート

既定は **`:8080`**（[`main.go`](main.go) の `addr`）。変更する場合はこの1箇所を編集してください。

### CSS の開発

- 一回ビルド: `npm run build`
- ウォッチ（保存のたびに `output.css` 更新）: `npm run watch`

HTML やクラスを増やしたあとは、ウォッチまたは `npm run build` で `output.css` を再生成してください。`input.css` の `@source` で `templates/**/*.html` と `node_modules/preline/dist/*.js` をスキャンしています。

## VS Code / Cursor（実行とデバッグ）

ワークスペースのルートを **`design-template1` フォルダ**にしたうえで利用してください（`${workspaceFolder}` がテンプレートと `static` を指す必要があります）。

| 構成名 | 内容 |
|--------|------|
| **design-template1: Go で起動** | Go のデバッグ起動のみ |
| **design-template1: CSS ビルド後に Go で起動** | `npm: build`（`tasks.json`）のあと Go 起動 |

Go 拡張機能（例: 公式 [Go](https://marketplace.visualstudio.com/items?itemName=golang.go)）が必要です。

## HTTP ルート

| メソッド・パス | 説明 |
|----------------|------|
| `GET /` | 管理画面デモ（`admin` テンプレート） |
| `GET /static/*` | 静的ファイル（CSS・JS） |
| その他 | `404` |

## スタイル・挙動のカスタムメモ

- **Preline 公式の Tailwind v4 手順**: `@import "tailwindcss"`、`@source` で Preline の JS、`variants.css`、`@plugin "@tailwindcss/forms"`（詳細は [Preline Quick setup](https://preline.co/docs/index.html)）
- **ツールチップの矢印**: `static/css/input.css` の `@layer components`（`[data-placement^="top"]` など）。色はライトで `#111827`、ダークで `#404040`（`prefers-color-scheme` と `.dark` 両対応）
- **バッジ**: 画面上のスタック表記・テーブルステータス・トップバーアバターは **ソリッド**（濃色背景＋白文字）で統一

## ライセンス・クレジット

- アプリケーションコード: このリポジトリの方針に従ってください。
- [Preline UI](https://preline.co/) は独自のライセンス条項があります。利用前に [Preline ライセンス](https://preline.co/docs/license.html) を確認してください。
