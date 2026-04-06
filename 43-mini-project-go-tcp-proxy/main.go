package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

func main() {
	// Command line flags
	listenAddr := flag.String("listen", ":8080", "Address to listen on")
	targetAddr := flag.String("target", "localhost:8000", "Target address to forward to")
	flag.Parse()

	fmt.Printf("🚀 Starting Go-TCP-Proxy\n")
	fmt.Printf("📍 Listening on: %s\n", *listenAddr)
	fmt.Printf("🎯 Forwarding to: %s\n", *targetAddr)
	fmt.Println("---------------------------------------------------------")

	// Start TCP listener
	listener, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("❌ Error listening on %s: %v", *listenAddr, err)
	}
	defer listener.Close()

	for {
		// Wait for incoming connections
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("⚠️  Error accepting connection: %v", err)
			continue
		}

		fmt.Printf("🤝 New connection from: %s\n", conn.RemoteAddr())
		go handleConnection(conn, *targetAddr)
	}
}

func handleConnection(source net.Conn, targetAddr string) {
	defer source.Close()

	// Dial the target address
	target, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("❌ Error dialing target %s: %v", targetAddr, err)
		return
	}
	defer target.Close()

	fmt.Printf("🔗 Connected to target: %s\n", targetAddr)

	var wg sync.WaitGroup
	wg.Add(2)

	// Copy from source to target (Proxying Requests)
	go func() {
		defer wg.Done()
		bytes, err := io.Copy(target, source)
		if err != nil {
			log.Printf("⚠️  Error copying from source to target: %v", err)
		}
		fmt.Printf("➡️  Sent %d bytes to target\n", bytes)
	}()

	// Copy from target to source (Proxying Responses)
	go func() {
		defer wg.Done()
		bytes, err := io.Copy(source, target)
		if err != nil {
			log.Printf("⚠️  Error copying from target to source: %v", err)
		}
		fmt.Printf("⬅️  Received %d bytes from target\n", bytes)
	}()

	wg.Wait()
	fmt.Printf("🛑 Connection closed for: %s\n", source.RemoteAddr())
}
