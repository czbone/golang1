# Go + Tailwind CSS サンプルプロジェクト

Go標準ライブラリとTailwind CSSを使用したシンプルなWebアプリケーションのサンプルです。

## 概要

- **バックエンド**: Go 1.21（標準ライブラリのみ）
- **フロントエンド**: Tailwind CSS v4.2.2
- **特徴**: 
  - 外部フレームワーク不要
  - 軽量かつ高速
  - ユーティリティファーストのスタイリング

## 必要な環境

- Go 1.21以上
- Node.js（npm）

## プロジェクト構造

```
web-tailwind/
├── main.go              # Goアプリケーションのエントリーポイント
├── go.mod               # Goモジュール定義
├── package.json         # npm依存関係とスクリプト
├── static/              # 静的ファイル
│   └── css/
│       ├── input.css    # Tailwind CSSソースファイル
│       └── output.css   # ビルドされたCSSファイル
└── templates/           # HTMLテンプレート
    └── index.html       # トップページ
```

## セットアップ

### 1. 依存関係のインストール

```bash
npm install
```

### 2. CSSのビルド

Tailwind CSSをビルドしてCSSファイルを生成します：

```bash
npm run build
```

開発中にCSSを自動的に再ビルドする場合は：

```bash
npm run watch
```

## 実行方法

### サーバーの起動

```bash
go run main.go
```

サーバーが起動すると、以下のメッセージが表示されます：

```
2026/04/04 14:38:41 Server starting on http://localhost:8080
```

### アクセス

ブラウザで以下のURLを開きます：

```
http://localhost:8080
```

### サーバーの停止

ターミナルで `Ctrl+C` を押してください。

## ビルド（バイナリ作成）

実行可能バイナリを作成する場合：

```bash
go build -o web-tailwind.exe main.go
```

作成されたバイナリを実行：

```bash
./web-tailwind.exe
```

## 開発ワークフロー

1. CSSを編集する場合は `npm run watch` を実行しておく
2. `main.go` やテンプレートを編集
3. サーバーを再起動して変更を確認

## カスタマイズ

- **ポート番号の変更**: [main.go](main.go#L37) の `addr := ":8080"` を編集
- **テンプレートの編集**: [templates/index.html](templates/index.html) を編集
- **スタイルの変更**: [static/css/input.css](static/css/input.css) を編集し、CSSを再ビルド

## トラブルシューティング

### ポートが既に使用されている場合

```
listen tcp :8080: bind: address already in use
```

別のアプリケーションがポート8080を使用している可能性があります。[main.go](main.go) のポート番号を変更してください。

### CSSが反映されない場合

`npm run build` を実行してCSSを再ビルドしてください。
