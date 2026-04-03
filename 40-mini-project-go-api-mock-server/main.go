package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// In-memory data store parsed from db.json
var (
	dataStore map[string][]map[string]interface{}
	mu        sync.RWMutex
)

// enableCORS middleware allows cross-origin requests
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// loadData reads the data from db.json and stores it in memory
func loadData(filename string) error {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	mu.Lock()
	defer mu.Unlock()
	
	err = json.Unmarshal(fileData, &dataStore)
	if err != nil {
		return fmt.Errorf("failed to parse JSON (must be an object of arrays): %w", err)
	}
	return nil
}

// requestHandler is the main router logic
func requestHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	if path == "" {
		handleHome(w, r)
		return
	}

	// Simple routing /collection or /collection/id
	parts := strings.Split(path, "/")
	
	mu.RLock()
	defer mu.RUnlock()

	collectionName := parts[0]
	collection, exists := dataStore[collectionName]
	
	if !exists {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Handled GET /collection
	if len(parts) == 1 {
		json.NewEncoder(w).Encode(collection)
		return
	}

	// Handled GET /collection/id
	if len(parts) == 2 {
		idParam := parts[1]
		
		for _, item := range collection {
			// Try to match id gracefully (could be string or float64 from JSON)
			idVal, hasId := item["id"]
			if hasId {
				strId := fmt.Sprintf("%v", idVal)
				if strId == idParam {
					json.NewEncoder(w).Encode(item)
					return
				}
			}
		}
		http.Error(w, "Item not found in collection", http.StatusNotFound)
		return
	}

	http.Error(w, "Bad request", http.StatusBadRequest)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `<h1>Go Mock API Server Running</h1>
	<p>Available endpoints based on db.json:</p><ul>`

	mu.RLock()
	defer mu.RUnlock()
	for key := range dataStore {
		html += fmt.Sprintf("<li><a href='/%s'>/%s</a></li>", key, key)
	}
	html += `</ul>`

	w.Write([]byte(html))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [path_to_db.json]\nDefaulting to './db.json'")
	}
	
	dbFile := "db.json"
	if len(os.Args) >= 2 {
		dbFile = os.Args[1]
	}

	fmt.Printf("Loading mock data from %s...\n", dbFile)
	err := loadData(dbFile)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	http.HandleFunc("/", enableCORS(requestHandler))

	port := "8080"
	fmt.Printf("🚀 Mock Server is running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
