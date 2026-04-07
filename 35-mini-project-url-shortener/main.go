package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

// --- Structs & State ---

type Shortener struct {
	mu    sync.RWMutex
	Links map[string]string `json:"links"`
}

const (
	port     = ":8080"
	dataFile = "links.json"
	idLength = 6
)

func main() {
	s := &Shortener{
		Links: make(map[string]string),
	}

	// 1. Load existing links from JSON
	s.load()

	// 2. Setup HTTP Routes
	http.HandleFunc("/shorten", s.handleShorten)
	http.HandleFunc("/", s.handleRedirect)

	fmt.Printf("🚀 URL Shortener starting on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// --- Handlers ---

func (s *Shortener) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Use POST to shorten URLs", http.StatusMethodNotAllowed)
		return
	}

	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	id := generateID(idLength)

	s.mu.Lock()
	s.Links[id] = url
	s.mu.Unlock()

	s.save()

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "✅ Success! Short URL: http://localhost%s/%s\n", port, id)
}

func (s *Shortener) handleRedirect(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:] // Extract ID from path
	if id == "" {
		fmt.Fprintln(w, "🐹 Welcome to Go-URL-Shortener! Use /shorten to condense your links.")
		return
	}

	s.mu.RLock()
	url, exists := s.Links[id]
	s.mu.RUnlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

// --- Helper Functions ---

func generateID(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, n)
	rand.Read(b)
	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b)
}

func (s *Shortener) save() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}
	os.WriteFile(dataFile, data, 0644)
}

func (s *Shortener) load() {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return // File doesn't exist yet, no problem
	}

	err = json.Unmarshal(data, s)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
	}
	fmt.Printf("📂 Loaded %d links from %s\n", len(s.Links), dataFile)
}
