# Day 31: Go-SSG - Static Site Generator 🛠️📄

## Overview

For Day 31, I built **Go-SSG**, a mini Static Site Generator. Go is the foundation of some of the world's fastest SSGs (like Hugo), and this project explores why. It transforms raw Markdown files into a complete, themed HTML website in seconds.

## Features

- **Markdown Transformation**: Converts `.md` files into valid HTML using the `gomarkdown` library.
- **Templating**: Uses Go's `html/template` to wrap content in a consistent layout.
- **Auto-Discovery**: Recursively walks the `content/` folder to build your site structure.
- **Asset Syncing**: Automatically copies CSS and other static files to the distribution folder.
- **Clean Builds**: Wipes the `public/` directory before every build to ensure no stale files remain.

## Project Structure

```text
31-mini-project-go-ssg/
├── content/        # Your Markdown source files
├── layout/         # HTML templates
├── static/         # CSS, Images, JS
├── public/         # THE BUILT WEBSITE (Generated)
└── main.go         # The SSG Engine
```

## How to Run

1. Navigate to the directory:
   ```bash
   cd 31-mini-project-go-ssg
   ```
2. Run the build:
   ```bash
   go run main.go
   ```
3. Open `public/index.html` in your browser to see your new static site!

## Learning Reflection

- **`html/template`**: Learned how to define blocks and pass dynamic content (like Page Titles and sanitized Body HTML) into a master layout.
- **FS Mirroring**: Practiced using `os.MkdirAll` and `filepath.Walk` to recreate a source directory structure in a destination folder.
- **Go Power**: Even with just a few hundred lines of code, Go provides the primitive tools to build professional-grade generation engines.

---
*Building the web, one Markdown file at a time. 🐹🛠️*
