# Day 35: Go-URL-Shortener - Persistent Web Service 🔗⚡

## Overview

For Day 35, I built **Go-URL-Shortener**, a high-performance web service that condenses long URLs into manageable 6-character IDs. This project focuses on building a stateful web application in Go, utilizing the `net/http` package for routing and `sync.RWMutex` for safe concurrent operations. It also features JSON-based persistence to ensure your shortened links survive a server restart.

## Features

- **Concurrent-Safe Storage**: Uses a thread-safe map protected by a Read-Write Mutex to handle simultaneous web requests.
- **RESTful API**: 
    - `POST /shorten`: Takes a long URL and generates a unique ID.
    - `GET /{id}`: Instantly redirects users to the destination (301 redirect).
- **Persistent Data**: Automatically saves the link registry to `links.json` on disk whenever a new link is created.
- **Auto-Recovery**: Loads existing mappings from disk when the server starts.
- **Fast Redirects**: Built for speed with O(1) lookup times for every redirection.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 35-mini-project-url-shortener
   ```
2. Run the server:
   ```bash
   go run main.go
   ```
3. Shorten a URL via `curl`:
   ```bash
   curl -X POST -d "url=https://golang.org" http://localhost:8080/shorten
   ```
4. Visit the returned short URL (e.g., `http://localhost:8080/abcdef`) in your browser to be redirected!

## Learning Reflection

- **Mutex vs. Concurrency**: Learned that while Go is great at multitasking, a raw `map` is not thread-safe. Using `sync.RWMutex` is essential for prevent crashes under heavy traffic.
- **State Serialization**: Used `encoding/json` to bridge the gap between in-memory speed and on-disk durability.
- **HTTP Middlewares**: Practiced basic HTTP request handling and status code management (301 for permanent redirects, 201 for creation).

---
*Bridging the web, one ID at a time. 🐹🔗*
