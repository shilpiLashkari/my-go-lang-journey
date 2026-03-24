package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// --- Structs ---

type client struct {
	conn     net.Conn
	nick     string
	outgoing chan string
}

type server struct {
	clients    map[string]client
	register   chan client
	unregister chan client
	broadcast  chan string
}

// --- Server Hub Logic ---

func (s *server) run() {
	for {
		select {
		case c := <-s.register:
			s.clients[c.nick] = c
			s.broadcast <- fmt.Sprintf("✨ %s joined the room!", c.nick)
			log.Printf("New connection: %s from %s", c.nick, c.conn.RemoteAddr())

		case c := <-s.unregister:
			delete(s.clients, c.nick)
			close(c.outgoing)
			s.broadcast <- fmt.Sprintf("👋 %s left the room.", c.nick)
			log.Printf("Disconnected: %s", c.nick)

		case msg := <-s.broadcast:
			for _, c := range s.clients {
				select {
				case c.outgoing <- msg:
				default:
					log.Printf("Warning: Dropping message for %s (buffer full)", c.nick)
				}
			}
		}
	}
}

// --- Client Handling ---

func handleConn(conn net.Conn, s *server) {
	defer conn.Close()

	// 1. Setup Client
	nick := strings.Split(conn.RemoteAddr().String(), ":")[1] // Default to port as nick
	c := client{
		conn:     conn,
		nick:     nick,
		outgoing: make(chan string, 10),
	}

	// 2. Start outgoing writer
	go func() {
		for msg := range c.outgoing {
			fmt.Fprintln(conn, msg)
		}
	}()

	s.register <- c
	fmt.Fprintln(conn, "🚀 Welcome to the Day 30 Milestone Chat!")
	fmt.Fprintln(conn, "Commands: /nick <name>, /list, /quit")
	fmt.Fprintln(conn, "----------------------------------------")

	// 3. Reader Loop
	input := bufio.NewScanner(conn)
	for input.Scan() {
		text := strings.TrimSpace(input.Text())
		if text == "" {
			continue
		}

		// Handle Commands
		if strings.HasPrefix(text, "/") {
			parts := strings.Split(text, " ")
			cmd := parts[0]

			switch cmd {
			case "/nick":
				if len(parts) > 1 {
					old := c.nick
					c.nick = parts[1]
					s.broadcast <- fmt.Sprintf("🆔 %s is now known as %s", old, c.nick)
				}
			case "/list":
				var nicks []string
				for n := range s.clients {
					nicks = append(nicks, n)
				}
				fmt.Fprintln(conn, "👥 Online: ", strings.Join(nicks, ", "))
			case "/quit":
				s.unregister <- c
				return
			default:
				fmt.Fprintln(conn, "❌ Unknown command.")
			}
			continue
		}

		// Broadcast public message
		s.broadcast <- fmt.Sprintf("💬 %s: %s", c.nick, text)
	}

	// Clean up on disconnect
	s.unregister <- c
}

// --- Main ---

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Fatal: Could not start server: %v", err)
	}

	s := &server{
		clients:    make(map[string]client),
		register:   make(chan client),
		unregister: make(chan client),
		broadcast:  make(chan string),
	}

	log.Println("🚀 Chat Server Milestone Day 30 — Starting on :8080")
	go s.run()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error: Connection failed: %v", err)
			continue
		}
		go handleConn(conn, s)
	}
}
