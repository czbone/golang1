package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// User はMongoDBのusersコレクションのドキュメント構造です
type User struct {
	ID    bson.ObjectID `bson:"_id"   json:"id"`
	Name  string        `bson:"name"  json:"name"`
	Email string        `bson:"email" json:"email"`
}

// BaseDb はMongoDBクライアントとデータベースをラップした構造体です
type BaseDb struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func main() {
	// MongoDBへ接続（コンテキストは接続確立のタイムアウト用）
	uri := "mongodb://mongo:mongo_password@localhost:27017/mongo_example?authSource=admin"
	//uri := "mongodb://mongo:mongo_password@localhost:27017/mongo_example"
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	// 接続確認（Ping）
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		log.Panicf("#failed to connect to MongoDB: %v", err)
	}
	log.Print("#connected to MongoDB")

	db := &BaseDb{
		Client: client,
		DB:     client.Database("sampledb"),
	}

	// データベース初期化（サンプルデータ挿入）
	initCtx, initCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer initCancel()
	if err := initializeDatabase(initCtx, db); err != nil {
		log.Panic(err)
	}

	r := gin.Default()

	// ルートエンドポイント
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to web-sqlx-mongodb API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"GET /":             "API information",
				"GET /ping":         "Health check with session test",
				"GET /users":        "Get all users",
				"POST /users":       "Create a user",
				"DELETE /users/:id": "Delete a user by ObjectID",
			},
		})
	})

	// ヘルスチェック（セッション＋トランザクションのテスト）
	// /ping        → 成功パターン（デフォルト）
	// /ping?error=true → エラーパターン（ロールバック）
	// ※ MongoDB のマルチドキュメントトランザクションはレプリカセット/シャードクラスタが必要
	r.GET("/ping", func(c *gin.Context) {
		shouldError := c.Query("error") == "true"

		err := db.DoInTx(c.Request.Context(), func(ctx context.Context) error {
			if shouldError {
				return errors.New("Error in session")
			}
			// セッション内でPingを実行してDB疎通を確認
			return db.Client.Ping(ctx, readpref.Primary())
		})

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "pong"})
	})

	// ユーザー一覧取得
	r.GET("/users", func(c *gin.Context) {
		var users []User
		cursor, err := db.DB.Collection("users").Find(c.Request.Context(), bson.D{})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(c.Request.Context())

		if err := cursor.All(c.Request.Context(), &users); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, users)
	})

	// ユーザー作成
	r.POST("/users", func(c *gin.Context) {
		var input struct {
			Name  string `json:"name"  binding:"required"`
			Email string `json:"email" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		user := User{
			ID:    bson.NewObjectID(),
			Name:  input.Name,
			Email: input.Email,
		}

		if _, err := db.DB.Collection("users").InsertOne(c.Request.Context(), user); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, user)
	})

	// ユーザー削除
	r.DELETE("/users/:id", func(c *gin.Context) {
		oid, err := bson.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid ObjectID format"})
			return
		}

		result, err := db.DB.Collection("users").DeleteOne(c.Request.Context(), bson.M{"_id": oid})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}

		c.JSON(200, gin.H{"message": "deleted"})
	})

	r.Run() // デフォルト: :8080
}

// DoInTx はMongoDBのセッションを開始してトランザクションを実行します。
// ※ マルチドキュメントトランザクションはレプリカセット/シャードクラスタ環境が必要です。
//
//	スタンドアロンの場合は StartTransaction がエラーになります。
func (db *BaseDb) DoInTx(ctx context.Context, f func(ctx context.Context) error) error {
	log.Print("#start session")

	session, err := db.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx context.Context) (any, error) {
		return nil, f(ctx)
	})
	if err != nil {
		log.Printf("#session error: %v", err)
		return err
	}

	log.Print("#commit session")
	return nil
}

// initializeDatabase はusersコレクションを確認し、空の場合はサンプルデータを挿入します
func initializeDatabase(ctx context.Context, db *BaseDb) error {
	collection := db.DB.Collection("users")

	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return err
	}
	log.Printf("#users collection document count: %d", count)

	// データが空の場合のみサンプルデータを挿入
	if count == 0 {
		sampleDocs := []interface{}{
			User{ID: bson.NewObjectID(), Name: "Alice", Email: "alice@example.com"},
			User{ID: bson.NewObjectID(), Name: "Bob", Email: "bob@example.com"},
			User{ID: bson.NewObjectID(), Name: "Charlie", Email: "charlie@example.com"},
		}

		result, err := collection.InsertMany(ctx, sampleDocs)
		if err != nil {
			return err
		}
		log.Printf("#inserted %d sample users", len(result.InsertedIDs))
	}

	return nil
}
