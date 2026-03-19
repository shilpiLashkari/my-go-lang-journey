package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const dataFile = "snippets.json"

// Snippet represents a single saved code snippet.
type Snippet struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Language  string    `json:"language"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// --- Persistence ---

func loadSnippets() []Snippet {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return []Snippet{}
	}
	var snippets []Snippet
	json.Unmarshal(data, &snippets)
	return snippets
}

func saveSnippets(snippets []Snippet) {
	data, err := json.MarshalIndent(snippets, "", "  ")
	if err != nil {
		fmt.Println("❌ Error saving snippets:", err)
		return
	}
	os.WriteFile(dataFile, data, 0644)
}

// --- Helpers ---

func nextID(snippets []Snippet) int {
	max := 0
	for _, s := range snippets {
		if s.ID > max {
			max = s.ID
		}
	}
	return max + 1
}

func printHeader() {
	fmt.Printf("%-5s %-8s %-25s %s\n", "ID", "LANG", "TITLE", "CREATED")
	fmt.Println(strings.Repeat("-", 60))
}

func printSnippetRow(s Snippet) {
	title := s.Title
	if len(title) > 24 {
		title = title[:21] + "..."
	}
	fmt.Printf("%-5d %-8s %-25s %s\n", s.ID, s.Language, title, s.CreatedAt.Format("2006-01-02"))
}

// --- Commands ---

func cmdAdd(args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	title := fs.String("title", "", "Snippet title (required)")
	lang := fs.String("lang", "text", "Programming language tag")
	content := fs.String("content", "", "Snippet content (required)")
	fs.Parse(args)

	if *title == "" || *content == "" {
		fmt.Println("❌ -title and -content are required.")
		return
	}

	snippets := loadSnippets()
	s := Snippet{
		ID:        nextID(snippets),
		Title:     *title,
		Language:  *lang,
		Content:   *content,
		CreatedAt: time.Now(),
	}
	snippets = append(snippets, s)
	saveSnippets(snippets)
	fmt.Printf("✅ Snippet #%d saved: '%s'\n", s.ID, s.Title)
}

func cmdList() {
	snippets := loadSnippets()
	if len(snippets) == 0 {
		fmt.Println("📭 No snippets saved yet. Use 'add' to create one.")
		return
	}
	fmt.Println("\n📋 YOUR SNIPPET LIBRARY")
	printHeader()
	for _, s := range snippets {
		printSnippetRow(s)
	}
}

func cmdSearch(args []string) {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	query := fs.String("query", "", "Search keyword")
	fs.Parse(args)

	if *query == "" {
		fmt.Println("❌ -query is required.")
		return
	}

	snippets := loadSnippets()
	q := strings.ToLower(*query)
	found := []Snippet{}
	for _, s := range snippets {
		if strings.Contains(strings.ToLower(s.Title), q) ||
			strings.Contains(strings.ToLower(s.Content), q) ||
			strings.Contains(strings.ToLower(s.Language), q) {
			found = append(found, s)
		}
	}

	if len(found) == 0 {
		fmt.Printf("🔍 No snippets match '%s'.\n", *query)
		return
	}

	fmt.Printf("\n🔍 Results for '%s':\n", *query)
	printHeader()
	for _, s := range found {
		printSnippetRow(s)
	}
}

func cmdShow(args []string) {
	fs := flag.NewFlagSet("show", flag.ExitOnError)
	id := fs.Int("id", 0, "Snippet ID to show")
	fs.Parse(args)

	if *id == 0 {
		fmt.Println("❌ -id is required.")
		return
	}

	snippets := loadSnippets()
	for _, s := range snippets {
		if s.ID == *id {
			fmt.Printf("\n📌 #%d: %s [%s]\n", s.ID, s.Title, s.Language)
			fmt.Println(strings.Repeat("-", 40))
			fmt.Println(s.Content)
			fmt.Println(strings.Repeat("-", 40))
			return
		}
	}
	fmt.Printf("❓ Snippet #%d not found.\n", *id)
}

func cmdDelete(args []string) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	id := fs.Int("id", 0, "Snippet ID to delete")
	fs.Parse(args)

	if *id == 0 {
		fmt.Println("❌ -id is required.")
		return
	}

	snippets := loadSnippets()
	newSnippets := []Snippet{}
	deleted := false
	for _, s := range snippets {
		if s.ID == *id {
			deleted = true
			continue
		}
		newSnippets = append(newSnippets, s)
	}

	if !deleted {
		fmt.Printf("❓ Snippet #%d not found.\n", *id)
		return
	}

	saveSnippets(newSnippets)
	fmt.Printf("🗑️  Snippet #%d deleted.\n", *id)
}

func printUsage() {
	fmt.Println("Usage: snippet-manager <command> [flags]")
	fmt.Println("Commands: add, list, search, show, delete")
}

// --- Main ---

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "add":
		cmdAdd(args)
	case "list":
		cmdList()
	case "search":
		cmdSearch(args)
	case "show":
		cmdShow(args)
	case "delete":
		cmdDelete(args)
	default:
		fmt.Printf("❌ Unknown command: %s\n", command)
		printUsage()
	}
}
