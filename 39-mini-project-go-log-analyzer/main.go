package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
)

// LogEntry represents a single line in a web server log (Combined Log Format)
type LogEntry struct {
	IP         string
	Timestamp  string
	Method     string
	Path       string
	StatusCode string
	Bytes      string
}

type Stats struct {
	TotalRequests   int
	StatusCodes     map[string]int
	IPCounts        map[string]int
	PathCounts      map[string]int
	TotalBytes      int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <log-file-path>")
		return
	}

	logFilePath := os.Args[1]
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Regex for Combined Log Format
	// Example: 127.0.0.1 - - [10/Oct/2000:13:55:36 -0700] "GET /index.html HTTP/1.1" 200 2326
	logRegex := regexp.MustCompile(`^(\S+) \S+ \S+ \[([^\]]+)\] "(\S+) (\S+) \S+" (\d{3}) (\d+|-)`)

	stats := Stats{
		StatusCodes: make(map[string]int),
		IPCounts:    make(map[string]int),
		PathCounts:  make(map[string]int),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := logRegex.FindStringSubmatch(line)

		if len(matches) < 7 {
			continue // Skip malformed lines
		}

		stats.TotalRequests++
		ip := matches[1]
		path := matches[4]
		status := matches[5]

		stats.IPCounts[ip]++
		stats.PathCounts[path]++
		stats.StatusCodes[status]++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning file: %v", err)
	}

	printSummary(stats)
}

func printSummary(stats Stats) {
	fmt.Println("========================================")
	fmt.Printf("🚀 Log Analysis Summary\n")
	fmt.Println("========================================")
	fmt.Printf("Total Requests: %d\n", stats.TotalRequests)
	fmt.Println("----------------------------------------")

	fmt.Println("✅ Status Codes:")
	for code, count := range stats.StatusCodes {
		fmt.Printf("  %s: %d\n", code, count)
	}
	fmt.Println("----------------------------------------")

	fmt.Println("📍 Top 5 IP Addresses:")
	topIPs := getTopN(stats.IPCounts, 5)
	for _, kv := range topIPs {
		fmt.Printf("  %-15s : %d\n", kv.Key, kv.Value)
	}
	fmt.Println("----------------------------------------")

	fmt.Println("🔗 Top 5 Requested Paths:")
	topPaths := getTopN(stats.PathCounts, 5)
	for _, kv := range topPaths {
		fmt.Printf("  %-30s : %d\n", kv.Key, kv.Value)
	}
	fmt.Println("========================================")
}

type KeyValue struct {
	Key   string
	Value int
}

func getTopN(m map[string]int, n int) []KeyValue {
	var kvList []KeyValue
	for k, v := range m {
		kvList = append(kvList, KeyValue{k, v})
	}

	sort.Slice(kvList, func(i, j int) bool {
		return kvList[i].Value > kvList[j].Value
	})

	if len(kvList) > n {
		return kvList[:n]
	}
	return kvList
}
