# sqlx-sqlite2

[sqlx](https://github.com/jmoiron/sqlx)、[modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)、[Gin](https://github.com/gin-gonic/gin) を使用した SQLite + Web API のサンプルプロジェクトです。

## 概要

このプロジェクトは以下の機能を実装しています。

- Pure Go SQLite ドライバ（CGO不要）を使用したデータベースアクセス
- Gin フレームワークによる Web サーバー
- トランザクション管理のカスタム実装
- データベースとテーブルの自動初期化
- コネクションプールの設定

## 使用技術

| パッケージ | バージョン | 説明 |
|---|---|---|
| Go | 1.26 | プログラミング言語 |
| [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin) | v1.12.0 | Web フレームワーク |
| [github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx) | v1.4.0 | `database/sql` の拡張ライブラリ |
| [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) | v1.46.1 | Pure Go 実装の SQLite3 ドライバ |

### modernc.org/sqlite について

`modernc.org/sqlite` は SQLite の C ソースコードを Pure Go に変換したライブラリです。  
`CGO_ENABLED=0` の環境（CGO が使えない Windows 環境など）でも動作します。

本プロジェクトでは `sql.Register` を使い `"sqlite3"` というドライバ名で登録しています。

```go
func init() {
    sql.Register("sqlite3", &sqlite.Driver{})
}
```

## 実行方法

### 通常の実行

```bash
go run main.go
```

### VS Code での実行とデバッグ

VS Code の「実行とデバッグ」パネル（`Ctrl+Shift+D`）から起動できます。

#### 前提条件

- [Go 拡張機能](https://marketplace.visualstudio.com/items?itemName=golang.go) のインストール

#### 起動方法

1. サイドバーの「実行とデバッグ」アイコンをクリックするか、`Ctrl+Shift+D` を押す
2. 上部のドロップダウンで **「Launch Package」** を選択
3. `F5` を押して実行（またはデバッグ）を開始する

ブレークポイントを設定してステップ実行することも可能です。

## API エンドポイント

サーバーは `http://localhost:8080` で起動します。

### GET /

API情報を返します。

```bash
curl http://localhost:8080/
```

```json
{
    "message": "Welcome to sqlx-sqlite2 API",
    "version": "1.0.0",
    "endpoints": {
        "GET /": "API information",
        "GET /ping": "Health check with transaction test"
    }
}
```

### GET /ping

トランザクションのテストエンドポイントです。クエリパラメータで成功/エラーパターンを切り替えられます。

**成功パターン（デフォルト）:**

```bash
curl http://localhost:8080/ping
```

```json
{
    "message": "pong"
}
```

**エラーパターン（`?error=true` を指定）:**

```bash
curl http://localhost:8080/ping?error=true
```

```json
{
    "error": "Error in transaction"
}
```

## プロジェクト構成

```
sqlx-sqlite2/
├── main.go          # メインプログラム
├── go.mod           # モジュール定義
├── go.sum           # 依存関係のチェックサム
├── sqlite/          # データベースファイル格納ディレクトリ
│   └── sample.db    # SQLite データベース（自動作成）
├── .vscode/
│   └── launch.json  # VS Code デバッグ設定
└── README.md        # このファイル
```

## 主な機能

### データベース初期化

初回起動時に自動的に以下を実行します。

- `users` テーブルの作成
- サンプルデータの挿入（Alice, Bob, Charlie）

### トランザクション管理

`BaseDb.DoInTx` メソッドでトランザクションを管理します。

```go
err := db.DoInTx(func(tx *sqlx.Tx) error {
    // トランザクション内の処理
    return nil
})
```

- エラーが返された場合は自動的にロールバック
- 正常終了時は自動コミット

### コネクションプール設定

```go
sqlxDb.SetMaxOpenConns(25)   // 最大オープン接続数
sqlxDb.SetMaxIdleConns(25)   // 最大アイドル接続数
sqlxDb.SetConnMaxLifetime(0) // 接続の最大ライフタイム（0は無制限）
```

## ビルド

```bash
go build -o sqlx-sqlite2.exe .
```

実行:

```bash
.\sqlx-sqlite2.exe
```

## データベーススキーマ

### users テーブル

| カラム | 型 | 説明 |
|---|---|---|
| id | INTEGER | 主キー（自動採番） |
| name | TEXT | ユーザー名 |
| email | TEXT | メールアドレス（ユニーク） |

## 開発のヒント

- `/ping` エンドポイントはデフォルトで成功レスポンスを返します
- `/ping?error=true` でエラーパターンをテストできます
- `sqlite/sample.db` を削除して再起動すると、データベースが再初期化されます
- トランザクション処理のカスタマイズは `BaseDb.DoInTx` メソッドを参照してください
