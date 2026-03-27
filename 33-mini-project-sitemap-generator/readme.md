# Day 33: Go-Sitemap-Generator - Web Crawler & SEO Mapper 🕷️🗺️

## Overview

For Day 33, I built **Go-Sitemap-Generator**, a specialized web crawler designed to map out an entire website's structure. It starts at a base URL, discovers internal links, and recursively traverses the domain to generate a standard `sitemap.xml` file. This is a powerful demonstration of Go's networking capabilities and data serialization.

## Features

- **Recursive Crawling**: Automatically explores all discoverable pages within the same domain.
- **Domain Guard**: Intelligently skips external links and social media redirects to keep the sitemap focused.
- **Visit Tracking**: Uses a `map` to prevent infinite loops and re-crawling of the same page.
- **XML Generation**: Leverages Go's `encoding/xml` to produce a valid, production-ready `sitemap.xml`.
- **Zero External Dependencies**: Built entirely using the Go standard library.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 33-mini-project-sitemap-generator
   ```
2. Run the crawler on a website (e.g., Go's official site):
   ```bash
   go run main.go -url https://go.dev -max 3
   ```
3. Check the produced `sitemap.xml` to see the results!

## Learning Reflection

- **Recursion vs. Iteration**: Explored how to use recursion for graph traversal while managing depth limits to prevent over-crawling.
- **XML Struct Tags**: Mastered using `xml:"urlset,attr"` and other tags to map Go structs to the specific schema required by search engines.
- **Resilient Parsing**: Implemented a more robust link extraction logic that handles relative paths (`/about`), absolute paths (`https://...`), and fragment identifiers (`#`).

---
*Mapping the web, one link at a time. 🐹🕷️*
