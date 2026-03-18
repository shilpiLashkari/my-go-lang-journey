# Day 24: Go-Repo-Scanner - Git Workspace Health Checker 📂🔍

## Overview

For Day 24, I built **Go-Repo-Scanner**, a developer utility that scans a directory for Git repositories and reports their status. As my workspace grows with multiple projects, I needed a quick way to see which repos are "dirty" (have uncommitted changes) without checking each folder manually. This project explores executing external CLI commands from within Go using `os/exec` and walking the file system.

## Features

- **Auto Discovery**: Recursively finds all Git repositories by detecting `.git` folders.
- **Branch Detection**: Reports the currently checked-out branch for every repo.
- **Dirty/Clean Status**: Instantly shows whether a repo has uncommitted changes.
- **Concurrent Scanning**: Checks all repos in parallel for maximum speed.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 24-mini-project-repo-scanner
   ```
2. Run from your parent workspace directory (e.g., scan `D:\My Workspace`):
   ```bash
   go run main.go -path "D:\My Workspace\Github Repo"
   ```

## Learning Reflection

- **`os/exec`**: Learned how to spawn child processes (like `git`) from Go and capture their `stdout`/`stderr` output as strings.
- **`filepath.Walk`**: Practiced traversing directory trees and using early-stopping logic to avoid descending into `.git` internal folders.
- **Practical Dev Tool**: Built something I will actually use every day to audit my own projects!

---
*Know the status of every repo at a glance! 🐹*
