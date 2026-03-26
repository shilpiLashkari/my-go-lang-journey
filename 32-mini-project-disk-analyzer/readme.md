# Day 32: Go-Disk-Analyzer - Visual Space Finder 📊📂

## Overview

For Day 32, I build **Go-Disk-Analyzer**, a command-line tool that visualizes disk usage. Unlike the standard `du` command which just lists numbers, this utility provides a colored bar chart showing the relative size of each subdirectory in a workspace. It’s perfect for finding hidden "space hogs" in your projects.

## Features

- **Recursive Scanning**: Accurately calculates the total size of entire folder structures.
- **Visual Bar Charts**: Renders colored ASCII bars for each folder based on its size relative to the largest folder.
- **Human-Readable Sizes**: Automatically formats sizes into B, KB, MB, GB, or TB.
- **Color-Coded Status**: Uses ANSI escape codes to make the output easy to scan.
- **Zero Dependencies**: Pure standard library implementation.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 32-mini-project-disk-analyzer
   ```
2. Run the analyzer on your current folder:
   ```bash
   go run main.go
   ```
3. Or scan a specific path:
   ```bash
   go run main.go -path "D:/My Projects"
   ```

## Learning Reflection

- **Recursive File Traversal**: Practiced walking directory trees while handling permission errors and symbolic links gracefully.
- **Data Scaling**: Learned how to normalize data (percentages) to a fixed width (20-character bars) for consistent UI layout.
- **String Formatting**: Mastered using `fmt.Printf` with custom padding and precision to create clean, tabular CLI reports.

---
*Finding space, one byte at a time. 🐹📊*
