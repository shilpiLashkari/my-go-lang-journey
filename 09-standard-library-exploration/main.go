package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Post represents a simple JSON structure from an API
type Post struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

func main() {
	fmt.Println("--- Day 9: Standard Library Exploration ---")

	// 1. Using 'os' and 'fmt'
	fmt.Printf("\nOS Package: Current Time: %s\n", time.Now().Format(time.RFC822))
	user := os.Getenv("USER")
	if user == "" {
		user = "Gopher"
	}
	fmt.Printf("Hello, %s! (Accessed via os.Getenv)\n", user)

	// 2. Using 'net/http' and 'io'
	fmt.Println("\nHTTP Package: Fetching a dummy post...")
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// 3. Using 'encoding/json'
	var post Post
	err = json.Unmarshal(body, &post)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	fmt.Printf("Post ID: %d\n", post.ID)
	fmt.Printf("Title: %s\n", post.Title)
	fmt.Printf("Body Snippet: %s...\n", post.Body[:30])

	fmt.Println("\nDay 9 Complete! The standard library is vast and powerful.")
}
