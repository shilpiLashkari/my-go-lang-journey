# Day 36: Go-Watchdog - Hot-Reload Development Utility 🐕🔄

## Overview

For Day 36, I built **Go-Watchdog**, a production-grade directory monitor. It watches your source code for any changes (saved, created, or deleted files) and automatically restarts your development command. Think of it as a language-agnostic "nodemon" or "Air" built entirely from scratch with the Go standard library. No more manual terminal restarts!

## Features

- **Recursive Monitoring**: Watches all subdirectories for changes, ensuring your entire project tree is covered.
- **Smart Debouncing**: Prevents multiple restarts during quick multi-file saves (e.g., when a formatter runs).
- **Process Orchestration**: Automatically kills the old process before spawning a new one, preventing port conflicts and memory leaks.
- **Polling Engine**: Uses an efficient file scanning strategy to detect modifications via `ModTime`.
- **Command Agnostic**: Use it to reload Go apps, Node.js scripts, Python servers, or even just run build scripts.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 36-mini-project-go-watchdog
   ```
2. **Watch and reload a simple Go app**:
   ```bash
   go run main.go -cmd "go run ../30-mini-project-tcp-chat-server/main.go"
   ```
3. **Trigger a reload**: Just save any file in the current directory or child folders.

## Learning Reflection

- **`os/exec` Orchestration**: Gained experience in managing child processes, including sending `SIGKILL` signals and handling input/output pipes.
- **Recursion vs. Performance**: Developed a recursive directory walker that builds a snapshot of the filesystem to compare against future states.
- **Polling Strategy**: Discovered how to use `time.Ticker` for a consistent monitoring cadence without taxing the CPU.

---
*Stay fast, stay reloading. 🐹🐕*
