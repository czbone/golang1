package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob(filepath.Join("templates", "*.html")))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := map[string]any{
		"Title": "Go + Tailwind CSS サンプル",
		"Cards": []map[string]string{
			{"Icon": "🚀", "Title": "高速", "Body": "Go の net/http は軽量かつ高速なWebサーバを提供します。"},
			{"Icon": "🎨", "Title": "モダンデザイン", "Body": "Tailwind CSS でユーティリティファーストのスタイリングを実現します。"},
			{"Icon": "🔧", "Title": "シンプル", "Body": "外部フレームワーク不要。標準ライブラリだけで動作します。"},
		},
	}
	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("template error: %v", err)
	}
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", indexHandler)

	addr := ":8080"
	log.Printf("Server starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
