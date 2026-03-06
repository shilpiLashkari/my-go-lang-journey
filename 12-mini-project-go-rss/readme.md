# Day 12: Go-RSS - Concurrent Feed Reader

## Overview

A CLI tool built with Go that concurrently fetches and parses RSS and Atom feeds from multiple sources. It leverages Go's powerful concurrency primitives (`goroutines` and `channels`) to provide a unified, sorted news timeline.

## Features

- **Concurrent Fetching**: Uses goroutines to request multiple feeds simultaneously.
- **Multi-Format Support**: Parses both RSS 2.0 and Atom feed formats using `encoding/xml`.
- **Unified Timeline**: Aggregates items from all sources and sorts them by publication date.
- **Configurable**: Easily add or remove feeds by editing the `feeds.json` file.

## Project Structure

- `main.go`: Contains the core logic for fetching, parsing, and displaying feeds.
- `feeds.json`: A configuration file where you can define your favorite news sources.

## How to Run

1. Navigate to this directory:
   ```bash
   cd 12-mini-project-go-rss
   ```
2. Build and run the application:
   ```bash
   go run main.go
   ```

## Learning Reflection

This project was a practical exercise in:

- Coordinating asynchronous tasks with `sync.WaitGroup`.
- Using `channels` for safe communication between goroutines.
- Manual XML unmarshaling for diverse data structures.
- Handling real-world networking and time-parsing challenges.
