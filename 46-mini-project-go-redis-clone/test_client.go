package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Printf("❌ Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	commands := []string{
		"SET user Shilpi",
		"GET user",
		"EXISTS user",
		"DEL user",
		"EXISTS user",
		"QUIT",
	}

	reader := bufio.NewReader(conn)
	for _, cmd := range commands {
		fmt.Printf("➡️  Sending: %s\n", cmd)
		fmt.Fprintf(conn, "%s\n", cmd)
		
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("❌ Error reading: %v\n", err)
			return
		}
		fmt.Printf("⬅️  Received: %s", response)
	}
}
