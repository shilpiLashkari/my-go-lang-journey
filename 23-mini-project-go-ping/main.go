package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// PingResult structures the data returned by our concurrent workers.
type PingResult struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Up         bool
	Error      string
}

func main() {
	fmt.Println("🚀 Starting Go-Ping (Concurrent Health Monitor)")
	fmt.Println("==================================================")

	// A sample list of websites to check
	websites := []string{
		"https://google.com",
		"https://github.com",
		"https://golang.org",
		"https://pkg.go.dev",
		"https://httpbin.org/delay/2", // This one intentionally takes 2 seconds
		"https://this-site-surely-does-not-exist-at-all.com", // DNS failure
	}

	// 1. Setup our Gathering Channel
	// We need a channel to funnel the scattered results back to the main thread.
	resultsChan := make(chan PingResult)

	// Record start time to prove concurrency
	overallStart := time.Now()

	// 2. Scatter the Work
	fmt.Printf("📡 Dispatching %d health checks concurrently...\n\n", len(websites))
	for _, url := range websites {
		// Launch a new goroutine for EVERY url immediately
		go checkSite(url, resultsChan)
	}

	// 3. Gather the Results
	// We know EXACTLY how many results to expect because we know the length of our slice.
	// This prevents the channel from deadlocking while waiting for infinite results.
	for i := 0; i < len(websites); i++ {
		res := <-resultsChan
		printResult(res)
	}

	// Print final stats to show this took way less time than a synchronous loop!
	duration := time.Since(overallStart)
	fmt.Println("==================================================")
	fmt.Printf("🏁 All checks complete in %v!\n", duration)
}

// checkSite takes a URL, performs a GET request with a strict timeout, and sends the result to the channel
func checkSite(url string, resultsChan chan<- PingResult) {
	// A custom HTTP client is crucial!
	// The default default http.Get() has NO timeout and could hang your goroutine forever.
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)

	// If a network-level error occurred (e.g., DNS failure, timeout)
	if err != nil {
		resultsChan <- PingResult{
			URL:      url,
			Duration: duration,
			Up:       false,
			Error:    err.Error(),
		}
		return
	}

	// ALWAYS close the body on successful response to prevent connection leaks
	defer resp.Body.Close()

	// Consider any 2xx or 3xx status code "UP"
	isUp := resp.StatusCode >= 200 && resp.StatusCode < 400

	// Send success data back to the main thread
	resultsChan <- PingResult{
		URL:        url,
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Up:         isUp,
	}
}

func printResult(res PingResult) {
	// Simple CLI formatting for readability
	statusStr := "❌ DOWN"
	if res.Up {
		statusStr = "✅  UP "
	}

	// Pad the URL for neat columns
	paddedURL := res.URL
	if len(paddedURL) > 40 {
		paddedURL = paddedURL[:37] + "..."
	}
	paddedURL = fmt.Sprintf("%-40s", paddedURL)

	if res.Error != "" {
		fmt.Printf("[%s] %s | ERROR: %s\n", statusStr, paddedURL, extractShortError(res.Error))
	} else {
		fmt.Printf("[%s] %s | HTTP %d | Time: %v\n", statusStr, paddedURL, res.StatusCode, res.Duration)
	}
}

// extractShortError tries to make Go's verbose networking errors slightly more readable for the CLI
func extractShortError(errStr string) string {
	if strings.Contains(errStr, "no such host") {
		return "DNS Lookup Failed"
	}
	if strings.Contains(errStr, "Client.Timeout") {
		return "Request Timed Out"
	}
	return "Connection Failed"
}
