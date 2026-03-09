# Day 15: Go-Watcher - Live Asset Monitor 🔍

## Overview

For Day 15, I've built **Go-Watcher**, a developer tool that monitors a directory for file changes and automatically runs a command. This is incredibly useful for high-productivity workflows where I want my code to re-run or my tests to trigger as soon as I save a file. This project taught me about polling the file system, handling child processes with `os/exec`, and using Go's `time` package for efficient loops.

## Features

- **Automated Execution**: Runs any shell command (like `go run main.go`) when a change is detected.
- **Efficient Polling**: Monitors modification times across large directories.
- **Customizable**:
  - Filter by file extensions (e.g., watch only `.go` or `.md`).
  - Set custom polling intervals.
  - Ignore specific directories (like `.git`).
- **Clean CLI**: Provides real-time feedback on what it's watching and when it's restarting.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 15-mini-project-go-watcher
   ```
2. Start watching and running a command (e.g., echo message):
   ```bash
   go run main.go -cmd "echo Change detected!" -interval 2s
   ```
3. Modify any `.go` file in the directory to see it in action!

## Learning Reflection

- **Process Control**: Learned how to start and properly manage child processes using `exec.Command`.
- **FS Mod-Times**: Understood how to use `os.Stat` and `filepath.Walk` to track changes without complex CGO dependencies.
- **Graceful Handling**: Implemented logic to ensure the previous command finishes or is logged before starting a new one.

---

_My Go development workflow just got a lot faster!_
