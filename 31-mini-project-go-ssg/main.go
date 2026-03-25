package main

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// PageData represents the variables passed to the HTML template.
type PageData struct {
	Title string
	Body  template.HTML
}

const (
	contentDir = "content"
	layoutFile = "layout/base.html"
	publicDir  = "public"
	staticDir  = "static"
)

func main() {
	fmt.Println("🚀 Starting Go-SSG Build Engine...")

	// 1. Clean public directory
	os.RemoveAll(publicDir)
	os.MkdirAll(publicDir, 0755)

	// 2. Load layout template
	tpl, err := template.ParseFiles(layoutFile)
	if err != nil {
		// Provide default files if they don't exist for the user yet
		setupProject()
		tpl, _ = template.ParseFiles(layoutFile)
	}

	// 3. Build content
	err = filepath.Walk(contentDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		if filepath.Ext(path) == ".md" {
			buildPage(path, tpl)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Build error: %v\n", err)
		return
	}

	// 4. Copy static assets
	copyStatic()

	fmt.Println("✨ Build Complete! Your site is in the 'public/' folder.")
}

func buildPage(srcPath string, tpl *template.Template) {
	// Read Markdown
	mdData, _ := os.ReadFile(srcPath)

	// Convert Markdown to HTML
	p := parser.NewWithExtensions(parser.CommonExtensions)
	mdHTML := markdown.ToHTML(mdData, p, html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags}))

	// Prepare output path
	relPath, _ := filepath.Rel(contentDir, srcPath)
	destPath := filepath.Join(publicDir, strings.TrimSuffix(relPath, ".md")+".html")
	os.MkdirAll(filepath.Dir(destPath), 0755)

	// Generate Page Meta
	title := strings.Title(strings.TrimSuffix(filepath.Base(srcPath), ".md"))
	data := PageData{
		Title: title,
		Body:  template.HTML(mdHTML),
	}

	// Execute Template and output file
	outFile, _ := os.Create(destPath)
	defer outFile.Close()
	tpl.Execute(outFile, data)

	fmt.Printf("📄 Generated: %s\n", destPath)
}

func copyStatic() {
	os.MkdirAll(filepath.Join(publicDir, staticDir), 0755)
	filepath.Walk(staticDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(staticDir, path)
		dest := filepath.Join(publicDir, staticDir, rel)
		
		srcFile, _ := os.Open(path)
		defer srcFile.Close()
		destFile, _ := os.Create(dest)
		defer destFile.Close()
		io.Copy(destFile, srcFile)
		return nil
	})
	fmt.Println("📁 Static assets copied.")
}

func setupProject() {
	fmt.Println("🏗️  Setting up default project structure...")
	os.MkdirAll(contentDir, 0755)
	os.MkdirAll("layout", 0755)
	os.MkdirAll(staticDir, 0755)

	// Sample Content
	os.WriteFile("content/index.md", []byte("# Welcome to Go-SSG\n\nThis site was built entirely in Go!"), 0644)
	os.WriteFile("content/about.md", []byte("# About This Journey\n\nDay 31 is all about Static Site Generation."), 0644)

	// Sample Layout
	layout := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.Title}} | My Go Journey</title>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <header><h1>Go-SSG Site</h1></header>
    <main>{{.Body}}</main>
    <footer>Built with Day 31 Go Project</footer>
</body>
</html>`
	os.WriteFile(layoutFile, []byte(layout), 0644)

	// Sample Static CSS
	css := "body { font-family: sans-serif; max-width: 800px; margin: 40px auto; padding: 20px; line-height: 1.6; background: #f4f4f9; }\nh1 { color: #00add8; }"
	os.WriteFile("static/style.css", []byte(css), 0644)
}
