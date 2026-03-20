# Day 26: Go-Notes-TUI - Colorful Interactive Terminal Notes 📝🎨

## Overview

For Day 26, I built **Go-Notes-TUI**, an interactive, colorful note-taking app that runs entirely in the terminal. While previous projects used sub-command routing, this one uses a live interactive menu loop — like a simple TUI (Text User Interface). It leverages ANSI escape codes directly to bring color and style to the terminal without any external dependencies.

## Features

- **Interactive Menu**: A colored menu that repeats until you choose to exit.
- **Create Notes**: Enter a title and multi-line content from the terminal.
- **List Notes**: See all notes in a formatted, colored table.
- **View Note**: Read the full content of any note.
- **Delete Note**: Remove a note by ID.
- **Persistent Storage**: Saves all notes to `notes.json`.
- **Zero External Dependencies**: Uses only Go's standard library — ANSI color codes built in.

## How to Run

```bash
cd 26-mini-project-notes-tui
go run main.go
```

## Learning Reflection

- **ANSI Escape Codes**: Learned to output colored, bold, and styled text using raw escape sequences like `\033[32m` without any library.
- **Interactive Loops**: Built a `for` loop menu that reads from `os.Stdin` using `bufio.Scanner`, a very different flow from sub-command CLI apps.
- **Input Handling**: Practiced reading multi-line user input safely in a terminal environment.

---
*The terminal is my canvas! 🐹*
