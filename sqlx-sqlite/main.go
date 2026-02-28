package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	sqlite "modernc.org/sqlite"
)

func init() {
	// modernc.org/sqlite (Pure Go) を "sqlite3" ドライバ名で登録する
	sql.Register("sqlite3", &sqlite.Driver{})
}

// User はユーザー情報を表す構造体
type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func main() {
	// インメモリ SQLite データベースを開く
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("データベースのオープンに失敗: %v", err)
	}
	defer db.Close()

	// テーブル作成
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id    INTEGER PRIMARY KEY AUTOINCREMENT,
			name  TEXT    NOT NULL,
			email TEXT    NOT NULL UNIQUE
		);
	`
	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("テーブル作成に失敗: %v", err)
	}
	fmt.Println("テーブルを作成しました")

	// レコード挿入
	users := []User{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
		{Name: "Charlie", Email: "charlie@example.com"},
	}
	for _, u := range users {
		_, err := db.Exec(
			"INSERT INTO users (name, email) VALUES (?, ?)",
			u.Name, u.Email,
		)
		if err != nil {
			log.Fatalf("レコード挿入に失敗: %v", err)
		}
	}
	fmt.Printf("%d 件のレコードを挿入しました\n", len(users))

	// 全レコード取得
	var result []User
	if err := db.Select(&result, "SELECT id, name, email FROM users ORDER BY id"); err != nil {
		log.Fatalf("レコード取得に失敗: %v", err)
	}

	fmt.Println("\n--- ユーザー一覧 ---")
	for _, u := range result {
		fmt.Printf("ID: %d  Name: %-10s  Email: %s\n", u.ID, u.Name, u.Email)
	}

	// 単一レコード取得 (sqlx.Get)
	var alice User
	if err := db.Get(&alice, "SELECT id, name, email FROM users WHERE name = ?", "Alice"); err != nil {
		log.Fatalf("単一レコード取得に失敗: %v", err)
	}
	fmt.Printf("\nGet で取得: %+v\n", alice)
}
