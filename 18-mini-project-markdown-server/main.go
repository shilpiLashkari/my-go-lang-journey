package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// Template for rendering HTML
const pageTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} | My Go Journey</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600&display=swap" rel="stylesheet">
    <style>
        body { font-family: 'Inter', sans-serif; line-height: 1.6; color: #333; max-width: 800px; margin: 40px auto; padding: 0 20px; background: #f9f9f9; }
        .container { background: white; padding: 40px; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        h1, h2, h3 { color: #00add8; }
        code { background: #eee; padding: 2px 5px; border-radius: 3px; font-family: monospace; }
        pre { background: #2d2d2d; color: #f8f8f2; padding: 15px; border-radius: 5px; overflow-x: auto; }
        a { color: #007bff; text-decoration: none; }
        a:hover { text-decoration: underline; }
        .back-link { display: block; margin-top: 20px; font-weight: 600; }
    </style>
</head>
<body>
    <div class="container">
        {{.Content}}
        <a href="/" class="back-link">🔙 Back to Home</a>
    </div>
</body>
</html>`

type PageData struct {
	Title   string
	Content template.HTML
}

func main() {
	http.HandleFunc("/", handler)
	port := ":8080"
	fmt.Printf("🚀 Markdown Server started at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		renderIndex(w)
		return
	}

	// Try to find the markdown file
	// I'll look back relative to the root of the project
	mdPath := filepath.Join("..", path)
	if !strings.HasSuffix(mdPath, ".md") {
		// Try appending readme.md if it's a directory
		mdPath = filepath.Join(mdPath, "readme.md")
	}

	content, err := os.ReadFile(mdPath)
	if err != nil {
		http.Error(w, "File not found: "+mdPath, http.StatusNotFound)
		return
	}

	// Convert Markdown to HTML
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	htmlRenderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags})
	
	htmlContent := markdown.ToHTML(content, p, htmlRenderer)

	// Fill template
	tmpl, _ := template.New("page").Parse(pageTemplate)
	data := PageData{
		Title:   filepath.Base(mdPath),
		Content: template.HTML(htmlContent),
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

func renderIndex(w http.ResponseWriter) {
	content := "# 📂 My Go Journey Documentation\n\nI can view all my progress below:\n\n"
	
	// Scan the project root
	entries, _ := os.ReadDir("..")
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			content += fmt.Sprintf("- [%s](/%s)\n", entry.Name(), entry.Name())
		}
	}
	content += "\n- [Main README.md](/README.md)"

	// Convert and serve index
	extensions := parser.CommonExtensions
	p := parser.NewWithExtensions(extensions)
	htmlContent := markdown.ToHTML([]byte(content), p, nil)

	tmpl, _ := template.New("page").Parse(pageTemplate)
	data := PageData{
		Title:   "Home",
		Content: template.HTML(htmlContent),
	}
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}
