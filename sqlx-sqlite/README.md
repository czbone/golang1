# sqlx-sqlite

[sqlx](https://github.com/jmoiron/sqlx) と [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) を使用して SQLite データベースにアクセスするサンプルプログラムです。

## 概要

このプロジェクトは以下を示すシンプルなサンプルです。

- インメモリ SQLite データベースの作成
- テーブルの作成
- レコードの挿入
- 複数レコードの取得 (`sqlx.Select`)
- 単一レコードの取得 (`sqlx.Get`)

## 使用技術

| パッケージ | バージョン | 説明 |
|---|---|---|
| Go | 1.26 | プログラミング言語 |
| [github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx) | v1.4.0 | `database/sql` の拡張ライブラリ |
| [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) | v1.46.1 | Pure Go 実装の SQLite3 ドライバ（CGO 不要） |

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

```bash
go run .
```

## 実行結果

```
テーブルを作成しました
3 件のレコードを挿入しました

--- ユーザー一覧 ---
ID: 1  Name: Alice       Email: alice@example.com
ID: 2  Name: Bob         Email: bob@example.com
ID: 3  Name: Charlie     Email: charlie@example.com

Get で取得: {ID:1 Name:Alice Email:alice@example.com}
```

## プロジェクト構成

```
sqlx-sqlite/
├── main.go        # メインプログラム
├── go.mod         # モジュール定義
├── go.sum         # 依存関係のチェックサム
├── .gitignore
└── .vscode/
    └── launch.json  # VS Code デバッグ設定
```
