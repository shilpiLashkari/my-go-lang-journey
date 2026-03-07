package main

import (
	"flag"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

func main() {
	// Define command-line flags
	host := flag.String("host", "scanme.nmap.org", "Host to scan")
	workers := flag.Int("workers", 100, "Number of concurrent workers")
	portRange := flag.String("ports", "1-1024", "Port range to scan (e.g. 1-100 or 80,443)")
	timeout := flag.Duration("timeout", 500*time.Millisecond, "Timeout for each port scan")
	flag.Parse()

	fmt.Printf("🔍 Starting Go-Scan on %s\n", *host)
	fmt.Printf("👷 Workers: %d | Timeout: %v\n", *workers, *timeout)
	fmt.Println("--------------------------------------------------")

	start := time.Now()

	var startPort, endPort int
	fmt.Sscanf(*portRange, "%d-%d", &startPort, &endPort)
	if startPort == 0 || endPort == 0 {
		startPort = 1
		endPort = 1024
	}

	ports := make(chan int, *workers)
	results := make(chan int)
	var openPorts []int

	// 1. Start the worker pool
	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go worker(*host, ports, results, *timeout, &wg)
	}

	// 2. Feed ports into the pool
	go func() {
		for i := startPort; i <= endPort; i++ {
			ports <- i
		}
		close(ports)
	}()

	// 3. Collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	for port := range results {
		if port != 0 {
			openPorts = append(openPorts, port)
			fmt.Printf("✅ Port %d is OPEN\n", port)
		}
	}

	sort.Ints(openPorts)

	duration := time.Since(start)
	fmt.Println("--------------------------------------------------")
	fmt.Printf("✨ Scan completed in %v\n", duration)
	fmt.Printf("📊 Found %d open ports.\n", len(openPorts))
}

func worker(host string, ports <-chan int, results chan<- int, timeout time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	for p := range ports {
		address := fmt.Sprintf("%s:%d", host, p)
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}
