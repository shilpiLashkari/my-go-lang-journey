# 11 - Mini Project - Go Journey Analytics 📊

For Day 11, I've built something unique and "meta"—a tool that analyzes my own progress in this repository! This project helps me see how much I've actually written and how my journey is structured.

## 🚀 Features

- **File Counter**: Counts how many Go programs and documentation files I've written.
- **Line Analytics**: Calculates the total lines of Go code across all days.
- **Concurrent Scanning**: Uses Goroutines to scan different directories at the same time for speed.
- **Progress Report**: Generates a clean dashboard summary of my learning.

## 🧠 Concepts Used

1. **File System walking**: Using `path/filepath` to navigate my project structure.
2. **Concurrency**: Using `sync.WaitGroup` to wait for multiple directory scanners to finish.
3. **Atomic Operations**: Using `sync/atomic` or `sync.Mutex` to safely update total counts from multiple goroutines.
4. **Formatting**: Using `text/tabwriter` or simple padding to create a nice-looking CLI dashboard.

## 🛠️ How to Use

1. Navigate to this directory: `cd 11-mini-project-journey-analytics`
2. Run the tool: `go run main.go`
3. Watch it analyze my entire Go journey!

---

_This tool gives me a clear view of my hard work so far!_
