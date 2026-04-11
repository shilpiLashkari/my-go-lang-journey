package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test_backends.go <port>")
		return
	}
	port := os.Args[1]
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from backend on port %s\n", port)
	})
	fmt.Printf("Starting mock backend on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
