package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// --- XML Schema ---

const (
	xmlNamespace = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

type URL struct {
	Loc string `xml:"loc"`
}

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

func main() {
	startURL := flag.String("url", "", "The website URL to crawl")
	maxPages := flag.Int("max", 20, "Maximum number of pages to crawl")
	flag.Parse()

	if *startURL == "" {
		fmt.Println("❌ Please provide a starting URL using -url")
		return
	}

	baseUrl, err := url.Parse(*startURL)
	if err != nil {
		fmt.Println("❌ Invalid URL:", err)
		return
	}

	fmt.Printf("🕷️  Starting crawl of: %s\n", baseUrl.String())

	// 1. Crawl
	pages := crawl(baseUrl.String(), *maxPages)

	// 2. Generate XML
	sitemap := URLSet{
		Xmlns: xmlNamespace,
		URLs:  make([]URL, len(pages)),
	}

	for i, p := range pages {
		sitemap.URLs[i] = URL{Loc: p}
	}

	// 3. Write to file
	fmt.Printf("📄 Generating sitemap with %d links...\n", len(pages))
	output, _ := xml.MarshalIndent(sitemap, "", "  ")
	header := []byte(xml.Header)
	os.WriteFile("sitemap.xml", append(header, output...), 0644)

	fmt.Println("✅ Done! Checkout sitemap.xml")
}

func crawl(baseUrl string, max int) []string {
	visited := make(map[string]bool)
	queue := []string{baseUrl}
	host := getHost(baseUrl)

	for len(queue) > 0 && len(visited) < max {
		curr := queue[0]
		queue = queue[1:]

		if visited[curr] {
			continue
		}

		fmt.Printf("   🔍 Crawling: %s\n", curr)
		links := getLinks(curr, baseUrl, host)
		visited[curr] = true

		for _, l := range links {
			if !visited[l] {
				queue = append(queue, l)
			}
		}
	}

	var results []string
	for k := range visited {
		results = append(results, k)
	}
	return results
}

func getLinks(target, base, host string) []string {
	resp, err := http.Get(target)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	// Poor man's link extraction using regex
	// In a real app, use golang.org/x/net/html
	re := regexp.MustCompile(`href=["']([^"']+)["']`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	var links []string
	for _, m := range matches {
		link := m[1]
		
		// Handle relative links
		u, err := url.Parse(link)
		if err != nil {
			continue
		}
		
		resolved := ""
		if u.IsAbs() {
			resolved = link
		} else {
			// Resolve relative path
			parent, _ := url.Parse(base)
			resolved = parent.ResolveReference(u).String()
		}

		// Only include links from the same host
		if getHost(resolved) == host {
			// Clean fragments
			clean := strings.Split(resolved, "#")[0]
			links = append(links, clean)
		}
	}

	return links
}

func getHost(u string) string {
	parsed, _ := url.Parse(u)
	return parsed.Host
}
