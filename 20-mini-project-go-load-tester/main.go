package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	targetURL   string
	numRequests int
	concurrency int
)

func init() {
	flag.StringVar(&targetURL, "url", "", "The target URL to load test (e.g., https://httpbin.org/get)")
	flag.IntVar(&numRequests, "requests", 100, "Total number of requests to send")
	flag.IntVar(&concurrency, "concurrency", 10, "Number of concurrent workers")
	flag.Parse()
}

// Stats tracking
var (
	successCount int64
	errorCount   int64
	totalTimeMs  int64 // Storing sum of all request durations in Ms
)

func main() {
	if targetURL == "" {
		fmt.Println("❌ Please provide a target URL using the -url flag.")
		return
	}

	fmt.Println("🚀 Go-Load-Tester Starting...")
	fmt.Printf("🎯 Target: %s\n", targetURL)
	fmt.Printf("📦 Total Requests: %d | ⚡ Concurrency: %d\n", numRequests, concurrency)
	fmt.Println("--------------------------------------------------")

	// Create a channel to act as the work queue
	jobs := make(chan int, numRequests)
	var wg sync.WaitGroup

	// Start the overall timer
	startTime := time.Now()

	// 1. Start the Workers
	for w := 1; w <= concurrency; w++ {
		wg.Add(1)
		go worker(&wg, jobs)
	}

	// 2. Load the Queue
	for i := 1; i <= numRequests; i++ {
		jobs <- i
	}
	close(jobs) // No more jobs to send

	// 3. Wait for all workers to finish
	wg.Wait()

	// Stop the overall timer
	duration := time.Since(startTime)

	// 4. Print the Report
	printReport(duration)
}

// worker constantly pulls jobs from the channel and fires HTTP requests
func worker(wg *sync.WaitGroup, jobs <-chan int) {
	defer wg.Done()

	// Reusable client with a timeout so workers don't get stuck forever
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	for range jobs {
		reqStart := time.Now()
		
		resp, err := client.Get(targetURL)
		
		reqDuration := time.Since(reqStart).Milliseconds()
		
		// Safely add this request's time to our global total
		atomic.AddInt64(&totalTimeMs, reqDuration)

		if err != nil {
			atomic.AddInt64(&errorCount, 1)
			continue
		}
		
		// We successfully got a response
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			atomic.AddInt64(&successCount, 1)
		} else {
			atomic.AddInt64(&errorCount, 1)
		}
		
		// Always close the body to prevent connection leaks!
		resp.Body.Close()
	}
}

func printReport(totalDuration time.Duration) {
	fmt.Println("\n📊 LOAD TEST REPORT")
	fmt.Println("==================================================")
	
	avgTime := int64(0)
	if successCount+errorCount > 0 {
		avgTime = totalTimeMs / (successCount + errorCount)
	}

	fmt.Printf("⏱️  Total Time Taken:     %v\n", totalDuration)
	fmt.Printf("✅ Successful Requests:  %d\n", successCount)
	fmt.Printf("🛑 Failed Requests:      %d\n", errorCount)
	fmt.Printf("📈 Average Req Time:     %d ms\n", avgTime)
	
	// Calculate requests per second
	rps := float64(numRequests) / totalDuration.Seconds()
	fmt.Printf("🏎️  Requests Per Second:  %.2f req/s\n", rps)
	fmt.Println("==================================================")
}
