package main

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var templates *template.Template

func loadTemplates() (*template.Template, error) {
	root := "templates"
	t := template.New("")
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".html" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = t.Parse(string(data))
		return err
	})
	return t, err
}

func init() {
	var err error
	templates, err = loadTemplates()
	if err != nil {
		log.Fatalf("templates: %v", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := map[string]any{
		"Title": "管理画面テンプレート",
		"Rows": []map[string]string{
			{"ID": "1", "Name": "山田 太郎", "Role": "管理者", "Status": "有効"},
			{"ID": "2", "Name": "佐藤 花子", "Role": "編集者", "Status": "有効"},
			{"ID": "3", "Name": "鈴木 一郎", "Role": "閲覧者", "Status": "停止"},
		},
	}
	if err := templates.ExecuteTemplate(w, "admin", data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("template error: %v", err)
	}
}

func main() {
	fsStatic := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fsStatic))
	http.HandleFunc("/", indexHandler)

	addr := ":8080"
	log.Printf("Server starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
