package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// --- Result Tracking ---

type Result struct {
	Latency    time.Duration
	StatusCode int
	Err        error
}

type Stats struct {
	mu           sync.Mutex
	SuccessCount int
	FailCount    int
	TotalLatency time.Duration
	MinLatency   time.Duration
	MaxLatency   time.Duration
	StatusCodes  map[int]int
}

func (s *Stats) Add(r Result) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if r.Err != nil {
		s.FailCount++
		return
	}

	s.SuccessCount++
	s.TotalLatency += r.Latency
	s.StatusCodes[r.StatusCode]++

	if s.MinLatency == 0 || r.Latency < s.MinLatency {
		s.MinLatency = r.Latency
	}
	if r.Latency > s.MaxLatency {
		s.MaxLatency = r.Latency
	}
}

func main() {
	url := flag.String("url", "", "The URL to benchmark")
	requests := flag.Int("requests", 100, "Total number of requests to perform")
	concurrency := flag.Int("concurrency", 10, "Number of multiple requests to make at a time")
	flag.Parse()

	if *url == "" {
		fmt.Println("❌ Please provide a URL with -url")
		return
	}

	stats := &Stats{
		StatusCodes: make(map[int]int),
	}

	fmt.Printf("⚡️ Benchmarking: %s\n", *url)
	fmt.Printf("📊 Requests: %d, Concurrency: %d\n\n", *requests, *concurrency)

	tasks := make(chan struct{}, *requests)
	results := make(chan Result, *requests)
	var wg sync.WaitGroup

	startTime := time.Now()

	// 1. Start Workers
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range tasks {
				reqStart := time.Now()
				resp, err := http.Get(*url)
				latency := time.Since(reqStart)

				res := Result{Latency: latency, Err: err}
				if resp != nil {
					res.StatusCode = resp.StatusCode
					resp.Body.Close()
				}
				results <- res
			}
		}()
	}

	// 2. Feed Tasks
	for i := 0; i < *requests; i++ {
		tasks <- struct{}{}
	}
	close(tasks)

	// 3. Collect Results
	go func() {
		for r := range results {
			stats.Add(r)
		}
	}()

	wg.Wait()
	close(results)

	totalTime := time.Since(startTime)

	// 4. Print Report
	printReport(stats, totalTime, *requests)
}

func printReport(s *Stats, total time.Duration, totalReqs int) {
	rps := float64(totalReqs) / total.Seconds()
	avgLatency := time.Duration(0)
	if s.SuccessCount > 0 {
		avgLatency = s.TotalLatency / time.Duration(s.SuccessCount)
	}

	fmt.Println("🏁 BENCHMARK COMPLETE")
	fmt.Println("========================================")
	fmt.Printf("Total Time:      %v\n", total.Round(time.Millisecond))
	fmt.Printf("Requests/sec:    %.2f\n", rps)
	fmt.Printf("Avg Latency:     %v\n", avgLatency.Round(time.Microsecond))
	fmt.Printf("Min Latency:     %v\n", s.MinLatency.Round(time.Microsecond))
	fmt.Printf("Max Latency:     %v\n", s.MaxLatency.Round(time.Microsecond))
	fmt.Println("----------------------------------------")
	fmt.Printf("Success:         %d\n", s.SuccessCount)
	fmt.Printf("Failed:          %d\n", s.FailCount)
	fmt.Println("Status Codes:")
	for code, count := range s.StatusCodes {
		fmt.Printf("  [%d]: %d\n", code, count)
	}
	fmt.Println("========================================")
}
