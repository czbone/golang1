package main

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var templates *template.Template

// Category represents a component category
type Category struct {
	ID   string
	Name string
}

// Section represents a section within a component page
type Section struct {
	ID    string
	Title string
}

// All component categories
var categories = []Category{
	{ID: "buttons", Name: "ボタン"},
	{ID: "forms", Name: "フォーム"},
	{ID: "cards", Name: "カード"},
	{ID: "alerts", Name: "アラート"},
	{ID: "modals", Name: "モーダル"},
	{ID: "badges", Name: "バッジ"},
	{ID: "tooltips", Name: "ツールチップ"},
	{ID: "progress", Name: "プログレス"},
	{ID: "dropdowns", Name: "ドロップダウン"},
	{ID: "tabs", Name: "タブ"},
	{ID: "breadcrumbs", Name: "パンくず"},
	{ID: "pagination", Name: "ページネーション"},
	{ID: "tables", Name: "テーブル"},
	{ID: "accordion", Name: "アコーディオン"},
}

// Sections for each category (for side navigation)
var categorySections = map[string][]Section{
	"buttons": {
		{ID: "default", Title: "基本ボタン"},
		{ID: "sizes", Title: "サイズ"},
		{ID: "outline", Title: "アウトライン"},
		{ID: "icons", Title: "アイコン付き"},
		{ID: "groups", Title: "ボタングループ"},
		{ID: "disabled", Title: "無効状態"},
	},
	"forms": {
		{ID: "inputs", Title: "テキスト入力"},
		{ID: "input-sizes", Title: "入力サイズ"},
		{ID: "textarea", Title: "テキストエリア"},
		{ID: "select", Title: "セレクト"},
		{ID: "checkbox", Title: "チェックボックス"},
		{ID: "radio", Title: "ラジオボタン"},
		{ID: "toggle", Title: "トグルスイッチ"},
		{ID: "file", Title: "ファイル入力"},
	},
	"cards": {
		{ID: "simple", Title: "シンプルカード"},
		{ID: "image", Title: "画像付きカード"},
		{ID: "horizontal", Title: "横並びカード"},
		{ID: "list", Title: "リスト付きカード"},
		{ID: "pricing", Title: "価格カード"},
	},
	"alerts": {
		{ID: "default", Title: "基本アラート"},
		{ID: "icons", Title: "アイコン付き"},
		{ID: "bordered", Title: "境界線あり"},
		{ID: "dismissible", Title: "閉じられる"},
		{ID: "additional", Title: "追加コンテンツ"},
	},
	"modals": {
		{ID: "default", Title: "基本モーダル"},
		{ID: "sizes", Title: "サイズ"},
		{ID: "form", Title: "フォームモーダル"},
		{ID: "confirmation", Title: "確認モーダル"},
	},
	"badges": {
		{ID: "default", Title: "基本バッジ"},
		{ID: "large", Title: "大きいバッジ"},
		{ID: "icons", Title: "アイコン付き"},
		{ID: "pill", Title: "ピル型"},
		{ID: "links", Title: "リンクバッジ"},
		{ID: "notification", Title: "通知バッジ"},
	},
	"tooltips": {
		{ID: "default", Title: "基本ツールチップ"},
		{ID: "positions", Title: "表示位置"},
		{ID: "styles", Title: "スタイル"},
	},
	"progress": {
		{ID: "default", Title: "基本プログレス"},
		{ID: "labels", Title: "ラベル付き"},
		{ID: "sizes", Title: "サイズ"},
		{ID: "spinner", Title: "ローディング"},
	},
	"dropdowns": {
		{ID: "default", Title: "基本ドロップダウン"},
		{ID: "divider", Title: "区切り線付き"},
		{ID: "header", Title: "ヘッダー付き"},
	},
	"tabs": {
		{ID: "default", Title: "基本タブ"},
		{ID: "pills", Title: "ピルタブ"},
	},
	"breadcrumbs": {
		{ID: "default", Title: "基本パンくず"},
		{ID: "solid", Title: "背景あり"},
	},
	"pagination": {
		{ID: "default", Title: "基本ページネーション"},
		{ID: "icons", Title: "アイコン付き"},
		{ID: "table", Title: "テーブル用ナビ"},
	},
	"tables": {
		{ID: "default", Title: "基本テーブル"},
		{ID: "striped", Title: "ストライプ"},
		{ID: "hover", Title: "ホバー"},
	},
	"accordion": {
		{ID: "default", Title: "基本アコーディオン"},
		{ID: "always-open", Title: "常に開く"},
	},
}

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
		"PageTitle":       "コンポーネント一覧",
		"CurrentCategory": "",
		"Categories":      categories,
	}
	if err := templates.ExecuteTemplate(w, "components-index", data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("template error: %v", err)
	}
}

func componentsIndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/components" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func componentCategoryHandler(category Category) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"PageTitle":       category.Name,
			"CurrentCategory": category.ID,
			"Categories":      categories,
			"Sections":        categorySections[category.ID],
		}
		templateName := "components-" + category.ID
		if err := templates.ExecuteTemplate(w, templateName, data); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("template error: %v", err)
		}
	}
}

func componentsRouterHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/components/")
	if path == "" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	categoryID := strings.Trim(path, "/")
	if categoryID != path {
		http.Redirect(w, r, "/components/"+categoryID, http.StatusMovedPermanently)
		return
	}

	for _, cat := range categories {
		if cat.ID == categoryID {
			componentCategoryHandler(cat)(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

func main() {
	fsStatic := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fsStatic))
	http.HandleFunc("/", indexHandler)

	// Component showcase routes
	http.HandleFunc("/components", componentsIndexHandler)
	http.HandleFunc("/components/", componentsRouterHandler)

	addr := ":8080"
	log.Printf("Server starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
