package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	// 1. Connect to the streamer
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("❌ Connection failed: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("🚀 Client connected and listening...")

	// 2. Append to the log file in a separate goroutine
	go func() {
		time.Sleep(1 * time.Second)
		f, err := os.OpenFile("app.log", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("❌ Error opening log file: %v\n", err)
			return
		}
		defer f.Close()

		messages := []string{"Hello from Go!", "Day 47 is awesome!", "Broadcast success!"}
		for _, msg := range messages {
			fmt.Printf("✍️  Writing to log: %s\n", msg)
			fmt.Fprintln(f, msg)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// 3. Read from the connection
	scanner := bufio.NewScanner(conn)
	timeout := time.After(5 * time.Second)
	
	for {
		select {
		case <-timeout:
			fmt.Println("⏲️  Test timeout reached.")
			return
		default:
			if scanner.Scan() {
				fmt.Printf("📥 Received: %s\n", scanner.Text())
			}
		}
	}
}
