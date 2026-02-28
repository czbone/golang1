package main

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	sqlite "modernc.org/sqlite"
)

func init() {
	// modernc.org/sqlite (Pure Go) を "sqlite3" ドライバ名で登録する
	sql.Register("sqlite3", &sqlite.Driver{})
}

type BaseDb struct {
	*sqlx.DB //匿名フィールド
}

type TranCallbackFunc func(*sqlx.Tx) error

func main() {

	// DBコネクションプールを作成（アプリケーション起動時に1回だけ）
	sqlxDb, err := sqlx.Connect("sqlite3", "./sqlite/sample.db")
	if err != nil {
		log.Panic(err)
	}
	defer sqlxDb.Close()

	// データベース初期化
	if err := initializeDatabase(sqlxDb); err != nil {
		log.Panic(err)
	}

	// コネクションプールの設定（オプション）
	sqlxDb.SetMaxOpenConns(25)   // 最大オープン接続数
	sqlxDb.SetMaxIdleConns(25)   // 最大アイドル接続数
	sqlxDb.SetConnMaxLifetime(0) // 接続の最大ライフタイム（0は無制限）

	db := &BaseDb{sqlxDb}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to sqlx-sqlite2 API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"GET /":     "API information",
				"GET /ping": "Health check with transaction test",
			},
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		// クエリパラメータで成功/エラーパターンを切り替え
		// /ping → 成功パターン（デフォルト）
		// /ping?error=true → エラーパターン
		shouldError := c.Query("error") == "true"

		// DBコネクションプールから接続を再利用
		err := db.DoInTx(func(tx *sqlx.Tx) error {
			if shouldError {
				return errors.New("Error in transaction")
			}
			return nil
		})

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run()
}

func checkErr(err error) {
	if err != nil {
		log.Print("#process exit by error")

		// 異常時は終了
		//log.Fatal(err)	// スタックトレースは出力しない =log.Fatal(err.Error())
		log.Panic(err) // スタックトレースも出力
	}
}

// func (db *BaseDb) DoInTx(f TranCallbackFunc(tx)) error {
func (db *BaseDb) DoInTx(f func(tx *sqlx.Tx) error) error {
	//func DoInTx(db *sqlx.DB, f func(tx *sqlx.Tx) (interface{}, error)) (interface{}, error) {

	log.Print("#start transaction")

	// トランザクション開始
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// DBデータ処理実行
	//err := f(tx) // ##### なぜ使えない? ##### NoNewVar occurs when a short variable declaration (':=') does not declare new variables. 「new」を使った値を返さないから?
	err = f(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// コミット
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	log.Print("#comit transaction")
	return nil
}

// initializeDatabase はデータベースのテーブルを初期化し、サンプルデータを挿入します
func initializeDatabase(db *sqlx.DB) error {
	// テーブル作成
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id    INTEGER PRIMARY KEY AUTOINCREMENT,
			name  TEXT    NOT NULL,
			email TEXT    NOT NULL UNIQUE
		);
	`
	if _, err := db.Exec(schema); err != nil {
		return err
	}
	log.Print("#database table initialized")

	// サンプルデータが既に存在するかチェック
	var count int
	if err := db.Get(&count, "SELECT COUNT(*) FROM users"); err != nil {
		return err
	}

	// データが空の場合のみサンプルデータを挿入
	if count == 0 {
		sampleUsers := []struct {
			Name  string
			Email string
		}{
			{Name: "Alice", Email: "alice@example.com"},
			{Name: "Bob", Email: "bob@example.com"},
			{Name: "Charlie", Email: "charlie@example.com"},
		}

		for _, u := range sampleUsers {
			_, err := db.Exec(
				"INSERT INTO users (name, email) VALUES (?, ?)",
				u.Name, u.Email,
			)
			if err != nil {
				return err
			}
		}
		log.Printf("#inserted %d sample users", len(sampleUsers))
	}

	return nil
}
