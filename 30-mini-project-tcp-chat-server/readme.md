# Day 30: Go-TCP-Chat-Server - Milestone: 30-Day Go Streak! 🚀🔥💬

## Overview

Welcome to the **Day 30 Milestone!** I've spent the last 30 days transitioning from MERN/Python into the world of high-performance Go. This project is the culmination of everything I've learned about concurrency, networking, and the "Go way" of building software.

**Go-TCP-Chat-Server** is a concurrent, real-time chat application. It allows multiple users to connect simultaneously via TCP (using tools like `telnet`), set nicknames, and chat in a global room. This project demonstrates the legendary Go concurrency model where dozens of clients are handled effortlessly by goroutines and channels.

## Features

- **Concurrent Hub**: A central server goroutine coordinates all message broadcasting and client management.
- **Dynamic Nicknames**: Users can set their identities with the `/nick` command.
- **Global Broadcast**: Messages are funneled through the server and instantly delivered to all active clients.
- **Server Commands**:
    - `/nick <name>` — Change your chat identity.
    - `/list` — See who else is currently online.
    - `/quit` — Gracefully disconnect from the server.
- **Live Notifications**: Displays when users join or leave the room.
- **Milestone Celebration**: Coded without external libraries, using pure Go standard networking.

## How to Run

1. **Start the Server**:
   ```bash
   cd 30-mini-project-tcp-chat-server
   go run main.go
   ```
2. **Connect from Terminals**:
   Open a new terminal and use `telnet` (or `nc` if on Linux/Mac):
   ```bash
   telnet localhost 8080
   ```
3. **Chat away!**

## Learning Reflection

- **Channel-based Concurrency**: Instead of using mutexes to lock a central user list, I used channels (`register`, `unregister`, `broadcast`) to communicate state changes to a single controlling goroutine. This avoids race conditions by design.
- **Goroutine per Connection**: Go handles 1 goroutine per socket perfectly. Each client has a dedicated reader that waits for incoming bytes and passes them to the hub.
- **Milestone Pride**: Reaching Day 30 feels incredible. Go has shifted the way I think about performance and system-level programming.

---
*30 days. 30 projects. 0 looking back. 🐹🚀*
