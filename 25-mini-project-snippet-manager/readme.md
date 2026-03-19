# Day 25: Go-Snippet-Manager - Personal Code Library 📋✨

## Overview

For Day 25, I built **Go-Snippet-Manager**, a terminal-based personal code library. Every developer accumulates useful snippets — handy one-liners, boilerplate setups, regex patterns — but they end up scattered across notes, Gists, or forgotten files. This tool gives me a searchable, persistent library managed entirely from the CLI, built in Go.

## Features

- **Add** snippets with a title, language tag, and content.
- **List** all saved snippets in a formatted table.
- **Search** snippets by keyword (matches title or content).
- **Show** the full content of any snippet by ID.
- **Delete** a snippet by ID.
- **Persists** everything to a local `snippets.json` file.

## How to Run

```bash
cd 25-mini-project-snippet-manager

# Add a snippet
go run main.go add -title "HTTP Server" -lang go -content "http.ListenAndServe(':8080', nil)"

# List all snippets
go run main.go list

# Search by keyword
go run main.go search -query "server"

# Show full content of snippet #1
go run main.go show -id 1

# Delete snippet #1
go run main.go delete -id 1
```

## Learning Reflection

- **Sub-command Routing**: Used `os.Args[1]` to dispatch to different command handlers, mirroring how tools like `git` or `docker` work.
- **JSON CRUD**: Practiced loading, mutating, and saving a slice of structs to a JSON file in a robust way.
- **Search Logic**: Used `strings.Contains` with `strings.ToLower` for case-insensitive search across multiple fields.

---
*My own personal Gist, powered by Go! 🐹*
