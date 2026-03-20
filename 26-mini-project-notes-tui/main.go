package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// --- ANSI Color Codes ---
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

const dataFile = "notes.json"

// Note represents a single note entry.
type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// --- Persistence ---

func loadNotes() []Note {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return []Note{}
	}
	var notes []Note
	json.Unmarshal(data, &notes)
	return notes
}

func saveNotes(notes []Note) {
	data, _ := json.MarshalIndent(notes, "", "  ")
	os.WriteFile(dataFile, data, 0644)
}

func nextID(notes []Note) int {
	max := 0
	for _, n := range notes {
		if n.ID > max {
			max = n.ID
		}
	}
	return max + 1
}

// --- UI Helpers ---

var scanner = bufio.NewScanner(os.Stdin)

func readLine(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func printBanner() {
	fmt.Println(colorCyan + colorBold)
	fmt.Println("  ╔══════════════════════════════╗")
	fmt.Println("  ║      📝  Go-Notes-TUI        ║")
	fmt.Println("  ╚══════════════════════════════╝")
	fmt.Println(colorReset)
}

func printMenu() {
	fmt.Println(colorYellow + "  Choose an option:" + colorReset)
	fmt.Println(colorGreen + "  [1]" + colorReset + " ✏️  Create Note")
	fmt.Println(colorGreen + "  [2]" + colorReset + " 📋 List All Notes")
	fmt.Println(colorGreen + "  [3]" + colorReset + " 👁  View Note")
	fmt.Println(colorGreen + "  [4]" + colorReset + " 🗑  Delete Note")
	fmt.Println(colorRed + "  [5]" + colorReset + " 🚪 Exit")
	fmt.Println()
}

// --- Commands ---

func createNote() {
	fmt.Println(colorCyan + "\n--- Create New Note ---" + colorReset)
	title := readLine("  Title: ")
	if title == "" {
		fmt.Println(colorRed + "  ❌ Title cannot be empty." + colorReset)
		return
	}
	content := readLine("  Content: ")

	notes := loadNotes()
	note := Note{
		ID:        nextID(notes),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}
	notes = append(notes, note)
	saveNotes(notes)
	fmt.Printf(colorGreen+"  ✅ Note #%d saved: '%s'\n"+colorReset, note.ID, note.Title)
}

func listNotes() {
	notes := loadNotes()
	fmt.Println(colorCyan + "\n--- Your Notes ---" + colorReset)
	if len(notes) == 0 {
		fmt.Println(colorYellow + "  📭 No notes yet. Create one!" + colorReset)
		return
	}
	fmt.Printf(colorBold+"  %-5s %-30s %s\n"+colorReset, "ID", "TITLE", "CREATED")
	fmt.Println("  " + strings.Repeat("-", 50))
	for _, n := range notes {
		title := n.Title
		if len(title) > 29 {
			title = title[:26] + "..."
		}
		fmt.Printf("  "+colorGreen+"%-5d"+colorReset+" %-30s %s\n", n.ID, title, n.CreatedAt.Format("2006-01-02 15:04"))
	}
}

func viewNote() {
	listNotes()
	input := readLine("\n  Enter Note ID to view: ")
	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println(colorRed + "  ❌ Invalid ID." + colorReset)
		return
	}
	for _, n := range loadNotes() {
		if n.ID == id {
			fmt.Printf(colorCyan+"\n  📌 #%d: %s\n"+colorReset, n.ID, n.Title)
			fmt.Println("  " + strings.Repeat("-", 40))
			fmt.Println("  " + n.Content)
			fmt.Println("  " + strings.Repeat("-", 40))
			return
		}
	}
	fmt.Println(colorRed + "  ❓ Note not found." + colorReset)
}

func deleteNote() {
	listNotes()
	input := readLine("\n  Enter Note ID to delete: ")
	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println(colorRed + "  ❌ Invalid ID." + colorReset)
		return
	}
	notes := loadNotes()
	newNotes := []Note{}
	deleted := false
	for _, n := range notes {
		if n.ID == id {
			deleted = true
			continue
		}
		newNotes = append(newNotes, n)
	}
	if !deleted {
		fmt.Println(colorRed + "  ❓ Note not found." + colorReset)
		return
	}
	saveNotes(newNotes)
	fmt.Printf(colorGreen+"  🗑  Note #%d deleted.\n"+colorReset, id)
}

// --- Main ---

func main() {
	for {
		clearScreen()
		printBanner()
		printMenu()

		choice := readLine("  > ")

		switch choice {
		case "1":
			createNote()
		case "2":
			listNotes()
		case "3":
			viewNote()
		case "4":
			deleteNote()
		case "5":
			fmt.Println(colorPurple + "\n  Goodbye! 👋\n" + colorReset)
			return
		default:
			fmt.Println(colorRed + "  ❌ Invalid choice. Try again." + colorReset)
		}

		readLine("\n  Press Enter to continue...")
	}
}
