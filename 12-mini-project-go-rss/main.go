package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// FeedConfig represents a feed entry in feeds.json
type FeedConfig struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Item represents a single news item from any feed type
type Item struct {
	Title       string
	Link        string
	Description string
	Published   time.Time
	Source      string
}

// RSS defines the structure for RSS 2.0 feeds
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// Atom defines the structure for Atom feeds
type Atom struct {
	XMLName xml.Name   `xml:"feed"`
	Title   string     `xml:"title"`
	Entries []AtomEntry `xml:"entry"`
}

type AtomEntry struct {
	Title     string    `xml:"title"`
	Link      AtomLink  `xml:"link"`
	Summary   string    `xml:"summary"`
	Published time.Time `xml:"published"`
	Updated   time.Time `xml:"updated"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
}

func main() {
	fmt.Println("🗞️  Go-RSS: Concurrent Feed Reader")
	fmt.Println("---------------------------------")

	// 1. Load feeds from configuration
	feeds, err := loadFeeds("feeds.json")
	if err != nil {
		fmt.Printf("❌ Error loading feeds: %v\n", err)
		return
	}

	// 2. Fetch feeds concurrently
	results := make(chan []Item)
	var wg sync.WaitGroup

	for _, feed := range feeds {
		wg.Add(1)
		go func(f FeedConfig) {
			defer wg.Done()
			items, err := fetchFeed(f)
			if err != nil {
				fmt.Printf("⚠️  Error fetching %s: %v\n", f.Name, err)
				return
			}
			results <- items
		}(feed)
	}

	// Close results channel once all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// 3. Collect and sort all items
	var allItems []Item
	for items := range results {
		allItems = append(allItems, items...)
	}

	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].Published.After(allItems[j].Published)
	})

	// 4. Display results
	if len(allItems) == 0 {
		fmt.Println("No items found.")
		return
	}

	fmt.Printf("\nFound %d latest news items:\n\n", len(allItems))
	for _, item := range allItems {
		fmt.Printf("[%s] %s\n", item.Source, item.Title)
		fmt.Printf("🔗 %s\n", item.Link)
		if item.Published.IsZero() {
			fmt.Printf("📅 Unknown date\n")
		} else {
			fmt.Printf("📅 %s\n", item.Published.Format("Jan 02, 2006 15:04"))
		}
		fmt.Println("---------------------------------")
	}
}

func loadFeeds(filename string) ([]FeedConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var feeds []FeedConfig
	err = json.Unmarshal(data, &feeds)
	return feeds, err
}

func fetchFeed(config FeedConfig) ([]Item, error) {
	resp, err := http.Get(config.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Try parsing as RSS
	var rss RSS
	if err := xml.Unmarshal(body, &rss); err == nil && len(rss.Channel.Items) > 0 {
		var items []Item
		for _, ri := range rss.Channel.Items {
			t, _ := parseTime(ri.PubDate)
			items = append(items, Item{
				Title:       ri.Title,
				Link:        ri.Link,
				Description: ri.Description,
				Published:   t,
				Source:      config.Name,
			})
		}
		return items, nil
	}

	// Try parsing as Atom
	var atom Atom
	if err := xml.Unmarshal(body, &atom); err == nil && len(atom.Entries) > 0 {
		var items []Item
		for _, ae := range atom.Entries {
			pub := ae.Published
			if pub.IsZero() {
				pub = ae.Updated
			}
			items = append(items, Item{
				Title:       ae.Title,
				Link:        ae.Link.Href,
				Description: ae.Summary,
				Published:   pub,
				Source:      config.Name,
			})
		}
		return items, nil
	}

	return nil, fmt.Errorf("unsupported or empty feed format")
}

func parseTime(timeStr string) (time.Time, error) {
	layouts := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		"Mon, 02 Jan 2006 15:04:05 MST", // Common but not standard RFC
		"2006-01-02T15:04:05Z",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, strings.TrimSpace(timeStr)); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unknown time format: %s", timeStr)
}
