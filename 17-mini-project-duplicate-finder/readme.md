# Day 17: Go-Duplicate-Finder - File Deduplication Tool 📂

## Overview

For Day 17, I've built **Go-Duplicate-Finder**, a utility to identify identical files within a directory tree. Unlike simple name checks, this tool uses SHA-256 cryptographic hashing to compare the actual content of the files. This ensures that even if two files have different names, they will be flagged as duplicates if their content is exactly the same. This project helped me practice deep file I/O, using the `crypto/sha256` package, and managing complex data with Go maps.

## Features

- **Content-Based Hashing**: Uses SHA-256 to ensure 100% accurate duplicate detection based on file content.
- **Recursive Scanning**: Uses `filepath.Walk` to scan all subdirectories.
- **Duplicate Grouping**: Automatically groups files with identical hashes for easy review.
- **Fast Hashing**: Reads files in chunks to handle large files efficiently.
- **CLI Interface**: Customizable scan directory via command-line flags.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 17-mini-project-duplicate-finder
   ```
2. Scan a specific directory (e.g., your Downloads or current folder):
   ```bash
   go run main.go -dir .
   ```

## Learning Reflection

- **Data Integrity**: Understood how unique hashes (SHA-256) acting as digital fingerprints can be used as keys in a map to identify duplicates.
- **Go Maps**: Practiced using a map where the value is a slice of strings (`map[string][]string`) to group multiple paths under a single content hash.
- **File Performance**: Used `io.Copy` with `sha256.New()` to stream file content into the hash function, which is much better for memory than reading the whole file at once.

---
*Clean files, clean mind! My repository is now more organized.*
