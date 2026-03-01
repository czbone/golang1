# web-sqlx-mongodb

[MongoDB Go Driver v2](https://www.mongodb.com/docs/drivers/go/current/) と [Gin](https://github.com/gin-gonic/gin) を使用した MongoDB + Web API のサンプルプロジェクトです。  
[web-sqlx-sqlite](../web-sqlx-sqlite) プロジェクトをベースに SQLite を MongoDB に置き換えたものです。

## 概要

このプロジェクトは以下の機能を実装しています。

- MongoDB Go Driver v2 を使用したデータベースアクセス
- Gin フレームワークによる Web サーバー
- セッション＋トランザクション管理のカスタム実装（`DoInTx`）
- コレクションの自動初期化とサンプルデータ挿入
- ユーザーの CRUD API エンドポイント

## 使用技術

| パッケージ | バージョン | 説明 |
|---|---|---|
| Go | 1.23 | プログラミング言語 |
| [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin) | v1.12.0 | Web フレームワーク |
| [go.mongodb.org/mongo-driver/v2](https://www.mongodb.com/docs/drivers/go/current/) | v2.2.0 | MongoDB 公式 Go ドライバ |

## 前提条件

MongoDB が起動していること（デフォルト: `localhost:27017`）。

### MongoDB のインストール例（Docker を使う場合）

```bash
docker run -d --name mongodb -p 27017:27017 mongo:latest
```

## 実行方法

```bash
go run main.go
```

## API エンドポイント

| メソッド | パス | 説明 |
|---|---|---|
| GET | `/` | API 情報 |
| GET | `/ping` | ヘルスチェック（セッションテスト） |
| GET | `/ping?error=true` | エラーパターン（セッションロールバック） |
| GET | `/users` | ユーザー一覧取得 |
| POST | `/users` | ユーザー作成 |
| DELETE | `/users/:id` | ユーザー削除（ObjectID 指定） |

### リクエスト例

#### ユーザー作成

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Dave", "email": "dave@example.com"}'
```

#### ユーザー一覧取得

```bash
curl http://localhost:8080/users
```

#### ユーザー削除

```bash
curl -X DELETE http://localhost:8080/users/<ObjectID>
```

## トランザクションについて

`DoInTx` はMongoDBのマルチドキュメントトランザクションを使用します。  
**注意**: マルチドキュメントトランザクションは **レプリカセット** または **シャードクラスタ** 環境が必要です。  
スタンドアロンの MongoDB ではトランザクションはエラーになります。  
ヘルスチェック（`/ping`）ではトランザクションの開始・コミット/ロールバックの動作を確認できます。

## データ構造

### User ドキュメント

```json
{
  "id": "ObjectID",
  "name": "string",
  "email": "string"
}
```

データは `sampledb` データベースの `users` コレクションに保存されます。  
起動時にコレクションが空の場合、Alice / Bob / Charlie の 3 件のサンプルデータが自動挿入されます。
