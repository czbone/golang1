# Copilot Instructions for webapp

## プロジェクト概要
- このリポジトリは Go + Gin + sqlx + SQLite + Tailwind CSS v4 で実装されたユーザー管理アプリです（Go 1.23）。
- エントリーポイントは `src/main.go`、ルーティングは `src/router.go` です。
- 画面は Go の `html/template` で `src/templates/*.html` を読み込みます。
- DB スキーマは `install/schema.sql`（`users` テーブル: id, name, email, created_at）、初期データは `install/seed.sql` です。

## 技術と主要ディレクトリ
- Web: Gin (`github.com/gin-gonic/gin`)
- DB: sqlx + modernc sqlite (`github.com/jmoiron/sqlx`, `modernc.org/sqlite`)
- CSS: Tailwind CSS v4 (`@tailwindcss/cli`) + `@tailwindcss/forms`
- UI コンポーネント: Preline UI v4 (`preline`)
- 設定: `src/config/config.go` — `Env{AppName, ServerPort, DatabasePath}` を返す `GetEnv()`
- コントローラ: `src/controllers/userController.go` — `UserController` 構造体
- DB クエリ層: `src/db/userDb.go` — `UserDb` 構造体（`*database.BaseDb` を埋め込み）
- DB 接続基盤: `src/lib/database/sqlite/sqlite.go` — `BaseDb` 構造体
- テンプレートヘルパー: `src/lib/template/data.go`（`MergeData`）, `template.go`（`LoadTemplates`）
- CSS 入力: `src/style/input.css`
- CSS 出力(生成物): `public/assets/css/style.css`
- Preline JS(生成物): `public/assets/js/preline.js`（`npm run build` 時に `node_modules/preline/dist/preline.js` からコピー）

## ルーティング（`src/router.go`）
| メソッド | パス | アクション |
|---------|------|-----------|
| GET | `/` | `/users` へリダイレクト (302) |
| GET | `/users` | `UserController.Index` |
| GET | `/users/new` | `UserController.New` |
| POST | `/users` | `UserController.Create` |
| GET | `/users/:id/edit` | `UserController.Edit` |
| POST | `/users/:id` | `UserController.Update` |
| POST | `/users/:id/delete` | `UserController.Delete` |

## DB 層の構造
- `BaseDb`（`src/lib/database/sqlite/sqlite.go`）が共通 DB 操作を提供する。
  - `GetDB()` — `sync.OnceValue` による単一接続（WAL モード + foreign_keys 有効）
  - `DoInTx(f func(tx *sqlx.Tx) error) error` — トランザクションヘルパー
  - `QueryRows(query, args...)` — 複数行を `[]map[string]interface{}` で返す
  - `QueryRow(query, args...)` — 単一行を `map[string]interface{}` で返す
  - `NewBaseDbWithDB(conn)` — テスト用の接続注入
- `UserDb`（`src/db/userDb.go`）は `*database.BaseDb` を埋め込み、ユーザ CRUD を実装する。

## 実装ルール
- 既存レイヤ構造を維持してください。
  - HTTP ハンドリングは `src/controllers/`
  - SQL/永続化ロジックは `src/db/`
  - DB 接続共通処理は `src/lib/database/sqlite/`
- 既存の命名規約とコメント言語（日本語コメント）を優先してください。
- DB の書き込み操作（INSERT / UPDATE / DELETE）は原則 `DoInTx` を使用してください。
- DB の読み取りには `QueryRows` / `QueryRow` ヘルパーを使用してください。
- テンプレートへ渡すデータは必ず `tmpl.MergeData(gin.H{...})` を使い、共通キー（`app_name`, `g_year`）を維持してください。
- 404 の描画は `tmpl.MergeData(gin.H{"page_title": "Not Found"})` を使い、`404.html` テンプレートを使用してください。
- バリデーションはコントローラ内の `validateUserInput` パターンに倣い、`errors []string` をテンプレートへ渡してください。
- 成功後のリダイレクトは `?flash=created|updated|deleted` クエリパラメータを付与するパターンに合わせてください。

## 変更時の注意点
- `public/assets/css/style.css` はビルド生成物です。直接編集せず `src/style/input.css` を編集してください。
- `public/assets/js/preline.js` はビルド生成物です。直接編集しないでください（`npm run build` で自動コピーされます）。
- SQLite の接続は `sync.OnceValue` で単一化されています。接続管理の方式を不用意に変更しないでください。
- 現在の環境変数名はコード実装を正とし、`SERVER_PORT`（デフォルト: `8080`）と `DATABASE_PATH`（デフォルト: `./user.sqlite3`）を使用してください。
- `AppName` は `config.go` で `"ユーザ管理"` にハードコードされています。
- 既存仕様を変える場合は、関連するテンプレート・ルート・DB クエリを一貫して更新してください。

## 変更後の確認コマンド
- Go のビルド確認: `go build ./src`
- CSS ビルド: `npm run build`
- 開発起動: `go run ./src`
- Makefile を使う場合: `make build`, `make run`, `make binary`

## Copilot への期待動作
- 変更は最小差分で行い、無関係なリファクタリングは避けてください。
- 新規機能追加時は、ルーティング (`src/router.go`)・コントローラ・DB 層・テンプレートの整合を必ず取ってください。
- SQL は既存スタイルに合わせてプレースホルダ `?` を使用してください。
- 可能であれば変更対象に応じて `go build ./src` または `npm run build` を実行し、結果を報告してください。
