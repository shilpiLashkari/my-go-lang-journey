package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// Broadcaster manages active TCP clients and routes log lines to them
type Broadcaster struct {
	mu      sync.Mutex
	clients map[net.Conn]chan string
	lines   chan string
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[net.Conn]chan string),
		lines:   make(chan string, 100),
	}
}

// Start initiates the broadcast loop
func (b *Broadcaster) Start() {
	for line := range b.lines {
		b.mu.Lock()
		for conn, ch := range b.clients {
			select {
			case ch <- line:
			default:
				log.Printf("⚠️  Client %s buffer full, skipping line", conn.RemoteAddr())
			}
		}
		b.mu.Unlock()
	}
}

// Register adds a new TCP connection to the broadcaster
func (b *Broadcaster) Register(conn net.Conn) {
	ch := make(chan string, 10)
	b.mu.Lock()
	b.clients[conn] = ch
	b.mu.Unlock()

	log.Printf("➕ Client connected: %s", conn.RemoteAddr())
	fmt.Fprintf(conn, "🚀 Connected to Go-Log-Streamer! Streaming logs...\n")

	// Send lines to client
	go func() {
		defer b.Unregister(conn)
		for line := range ch {
			_, err := fmt.Fprint(conn, line)
			if err != nil {
				return // Disconnected
			}
		}
	}()
}

// Unregister removes a connection
func (b *Broadcaster) Unregister(conn net.Conn) {
	b.mu.Lock()
	if ch, ok := b.clients[conn]; ok {
		close(ch)
		delete(b.clients, conn)
		conn.Close()
		log.Printf("➖ Client disconnected: %s", conn.RemoteAddr())
	}
	b.mu.Unlock()
}

// TailFile watches a file for new line appends
func TailFile(filename string, lines chan<- string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("❌ Error opening file: %v", err)
	}
	defer f.Close()

	// Move to the end of the file initially
	_, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		log.Fatalf("❌ Error seeking end of file: %v", err)
	}

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			// Wait for more data
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			log.Printf("⚠️  Error reading file: %v", err)
			break
		}
		lines <- line
	}
}

func main() {
	filePath := flag.String("file", "app.log", "The log file to stream")
	port := flag.String("port", "8080", "TCP port to listen on")
	flag.Parse()

	// 1. Ensure file exists or create it
	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		log.Printf("📝 File %s does not exist, creating it...", *filePath)
		os.WriteFile(*filePath, []byte(""), 0644)
	}

	broadcaster := NewBroadcaster()

	// 2. Start Broadcaster Loop
	go broadcaster.Start()

	// 3. Start Tailing the file in background
	go TailFile(*filePath, broadcaster.lines)

	// 4. Start TCP Server
	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("❌ Failed to start TCP server: %v", err)
	}
	defer listener.Close()

	fmt.Printf("🚀 Go-Log-Streamer live on port %s\n", *port)
	fmt.Printf("📁 Monitoring: %s\n", *filePath)
	fmt.Println("Connect with: nc localhost", *port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("⚠️  Accept error: %v", err)
			continue
		}
		broadcaster.Register(conn)
	}
}
