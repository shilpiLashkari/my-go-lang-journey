package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	filePath string
	workers  int
	timeout  time.Duration
)

func init() {
	flag.StringVar(&filePath, "file", "links.txt", "Path to text file containing URLs")
	flag.IntVar(&workers, "workers", 10, "Number of concurrent workers")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "HTTP request timeout")
	flag.Parse()
}

type Result struct {
	URL    string
	Status int
	Alive  bool
	Error  error
}

func main() {
	fmt.Println("🔗 Go-Link-Checker: Starting...")
	fmt.Printf("📂 File: %s | 👷 Workers: %d | ⏱️ Timeout: %v\n", filePath, workers, timeout)

	links, err := readLinks(filePath)
	if err != nil {
		fmt.Printf("❌ Error reading file: %v\n", err)
		return
	}

	if len(links) == 0 {
		fmt.Println("❓ No links found in the file.")
		return
	}

	jobs := make(chan string, len(links))
	results := make(chan Result, len(links))
	var wg sync.WaitGroup

	// Start workers
	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go worker(&wg, jobs, results)
	}

	// Send jobs
	for _, link := range links {
		jobs <- link
	}
	close(jobs)

	// Wait for workers to finish in a separate goroutine
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and report results
	aliveCount := 0
	brokenCount := 0

	fmt.Println("\n🔎 Verification Progress:")
	fmt.Println("---------------------------")
	for res := range results {
		if res.Alive {
			fmt.Printf("✅ [%d] %s\n", res.Status, res.URL)
			aliveCount++
		} else {
			if res.Error != nil {
				fmt.Printf("❌ [ERR] %s - %v\n", res.URL, res.Error)
			} else {
				fmt.Printf("❌ [%d] %s\n", res.Status, res.URL)
			}
			brokenCount++
		}
	}

	fmt.Println("\n📊 FINAL SUMMARY")
	fmt.Println("========================")
	fmt.Printf("🌟 Alive Links:  %d\n", aliveCount)
	fmt.Printf("💔 Broken Links: %d\n", brokenCount)
	fmt.Printf("🏁 Total Links:  %d\n", len(links))
	fmt.Println("========================")
}

func readLinks(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var links []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		link := scanner.Text()
		if link != "" {
			links = append(links, link)
		}
	}
	return links, scanner.Err()
}

func worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- Result) {
	defer wg.Done()

	client := http.Client{
		Timeout: timeout,
	}

	for url := range jobs {
		resp, err := client.Head(url)
		if err != nil {
			results <- Result{URL: url, Alive: false, Error: err}
			continue
		}
		
		alive := resp.StatusCode >= 200 && resp.StatusCode < 400
		results <- Result{URL: url, Status: resp.StatusCode, Alive: alive}
		resp.Body.Close()
	}
}
