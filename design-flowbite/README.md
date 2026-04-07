# design-flowbite

Go 標準ライブラリで動作する、Flowbite コンポーネントカタログです。  
バックエンドは Go (`net/http`, `html/template`) のみ、フロントは Tailwind CSS v4 + Flowbite v4 で構成しています。

## 技術スタック

| 領域 | 内容 |
|------|------|
| バックエンド | Go 1.21+、`net/http`、`html/template` |
| スタイル | Tailwind CSS v4、`@tailwindcss/cli` |
| UI コンポーネント | [Flowbite v4](https://flowbite.com/) (`flowbite` npm パッケージ) |
| フォーム補助 | `@tailwindcss/forms` |

Flowbite のインタラクション（ドロップダウン、モーダル、ツールチップなど）は `static/js/flowbite.min.js` で有効化されます。  
このファイルは `npm run build` 時に `scripts/copy-flowbite.js` で `node_modules` からコピーされます。

## 画面仕様

### 1. コンポーネント一覧 (`GET /`)

- Hero + カテゴリカードの一覧画面
- 各カードからカテゴリ別デモページへ遷移
- ヘッダーのカテゴリナビ（デスクトップ）とドロワー（モバイル）

### 2. カテゴリ別コンポーネントページ (`GET /components/{category}`)

- 左側（LG 以上）にページ内セクションナビ
- 各カテゴリごとに複数の実装例を掲載
- `#section-id` でページ内ジャンプ可能

対応カテゴリ:

- `buttons`
- `forms`
- `cards`
- `alerts`
- `modals`
- `badges`
- `tooltips`
- `progress`
- `dropdowns`
- `tabs`
- `breadcrumbs`
- `pagination`
- `tables`
- `accordion`

テンプレートは `templates/**/*.html` を再帰読み込みし、`ExecuteTemplate` で以下の名前を描画します。

- 一覧: `components-index`
- カテゴリ: `components-{category}`（例: `components-buttons`）

## ディレクトリ構成

```
design-flowbite/
├── main.go
├── go.mod
├── package.json
├── scripts/
│   ├── copy-flowbite.js
│   └── translate-components-ja.ps1
├── static/
│   ├── css/
│   │   ├── input.css
│   │   └── output.css
│   └── js/
│       └── flowbite.min.js
├── templates/
│   ├── components/
│   │   ├── index.html
│   │   ├── buttons.html
│   │   ├── forms.html
│   │   └── ...
│   └── partials/
│       ├── components-header.html
│       └── components-sidenav.html
└── README.md
```

## 必要な環境

- Go 1.21 以上
- Node.js / npm

## セットアップと起動

```bash
cd design-flowbite
npm install
npm run build
go run .
```

ブラウザで `http://localhost:8080/` を開きます。

## npm スクリプト

- `npm run build:css`: Tailwind をビルドして `static/css/output.css` を生成
- `npm run copy:flowbite`: `flowbite.min.js` を `static/js/` へコピー
- `npm run build`: CSS ビルド + Flowbite JS コピー
- `npm run watch`: CSS をウォッチビルド

`static/css/input.css` では `@source "../../templates/**/*.html"` と `@source "../../node_modules/flowbite"` を指定してクラス抽出しています。

## HTTP ルート

| メソッド・パス | 説明 |
|----------------|------|
| `GET /` | コンポーネント一覧 (`components-index`) |
| `GET /components` | `/` へリダイレクト |
| `GET /components/{category}` | カテゴリページ (`components-{category}`) |
| `GET /static/*` | 静的ファイル配信 |
| その他 | `404` |

## ポート

既定は `:8080`（`main.go` の `addr`）。必要に応じて変更してください。

## ライセンス・クレジット

- アプリケーションコード: このリポジトリの方針に従ってください。
- Flowbite の利用条件は公式ドキュメントを確認してください。
