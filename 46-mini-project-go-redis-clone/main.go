package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

// Storage handles the in-memory key-value map with a mutex for thread safety
type Storage struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

// SET saves a key-value pair
func (s *Storage) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// GET retrieves a value by key
func (s *Storage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

// DEL deletes a key
func (s *Storage) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[key]
	if ok {
		delete(s.data, key)
	}
	return ok
}

// EXISTS checks if a key exists
func (s *Storage) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[key]
	return ok
}

func handleConnection(conn net.Conn, storage *Storage) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	fmt.Printf("⚡ New connection from %s\n", conn.RemoteAddr())

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])

		switch command {
		case "SET":
			if len(parts) < 3 {
				conn.Write([]byte("ERR missing arguments: SET <key> <value>\n"))
				continue
			}
			key := parts[1]
			value := strings.Join(parts[2:], " ")
			storage.Set(key, value)
			conn.Write([]byte("OK\n"))

		case "GET":
			if len(parts) < 2 {
				conn.Write([]byte("ERR missing arguments: GET <key>\n"))
				continue
			}
			key := parts[1]
			val, ok := storage.Get(key)
			if !ok {
				conn.Write([]byte("(nil)\n"))
			} else {
				conn.Write([]byte(val + "\n"))
			}

		case "DEL":
			if len(parts) < 2 {
				conn.Write([]byte("ERR missing arguments: DEL <key>\n"))
				continue
			}
			key := parts[1]
			removed := storage.Del(key)
			if removed {
				conn.Write([]byte("(integer) 1\n"))
			} else {
				conn.Write([]byte("(integer) 0\n"))
			}

		case "EXISTS":
			if len(parts) < 2 {
				conn.Write([]byte("ERR missing arguments: EXISTS <key>\n"))
				continue
			}
			key := parts[1]
			exists := storage.Exists(key)
			if exists {
				conn.Write([]byte("(integer) 1\n"))
			} else {
				conn.Write([]byte("(integer) 0\n"))
			}

		case "QUIT":
			conn.Write([]byte("Bye!\n"))
			return

		default:
			conn.Write([]byte(fmt.Sprintf("ERR unknown command: '%s'\n", command)))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}
	fmt.Printf(" 👋 Connection closed from %s\n", conn.RemoteAddr())
}

func main() {
	storage := NewStorage()
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	fmt.Println("🚀 Go-Redis-Clone Lite started on :6379")
	fmt.Println("Commands: SET, GET, DEL, EXISTS, QUIT")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept failed: %v", err)
			continue
		}
		go handleConnection(conn, storage)
	}
}
