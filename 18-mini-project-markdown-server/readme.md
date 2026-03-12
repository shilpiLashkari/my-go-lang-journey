# Day 18: Go-Markdown-Server - Local Docs Viewer 🌐

## Overview

For Day 18, I've built **Go-Markdown-Server**, a personal documentation viewer that turns this repository into a live website. After exploring file hashing and system monitors, I wanted to see how Go handles web serving and external libraries. This tool starts a local server that renders my `.md` files as styled HTML, making it much easier for me to review my journey in the browser.

## Features

- **Live Rendering**: Converted Markdown to clean HTML on-the-fly.
- **Dynamic Content**: Serves any `.md` file in the project.
- **Styled Layout**: Uses embedded CSS to give my notes a premium feel.
- **Go Modules**: Managed my first major external dependency using `go.mod`.
- **Directory Discovery**: Automatically finds and lists files for me.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 18-mini-project-markdown-server
   ```
2. Install the rendering engine:
   ```bash
   go get github.com/gomarkdown/markdown
   ```
3. Run the server:
   ```bash
   go run main.go
   ```
4. Open your browser and visit: `http://localhost:8080`

## Learning Reflection

- **HTTP Serving**: Deepened my knowledge of `net/http` and how to handle file routing.
- **Dependency Management**: Learned how to effectively use `go get` and how `go.mod` tracks third-party rendering libraries.
- **HTML Templates**: Used `html/template` to safely inject my rendered markdown content into a beautiful web shell.
- **Media Handling**: Understood the importance of setting the correct Content-Type for HTML responses.

---
*Now I can browse my Go journey like a pro!*
