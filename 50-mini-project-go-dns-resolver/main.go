package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// ANSI color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
)

// DNSResult holds the outcome of a single DNS lookup.
type DNSResult struct {
	Domain     string
	RecordType string
	Results    []string
	Duration   time.Duration
	Error      error
}

// lookupA resolves A/AAAA records (IP addresses) for a domain.
func lookupA(domain string) DNSResult {
	start := time.Now()
	addrs, err := net.LookupHost(domain)
	dur := time.Since(start)

	return DNSResult{
		Domain:     domain,
		RecordType: "A/AAAA",
		Results:    addrs,
		Duration:   dur,
		Error:      err,
	}
}

// lookupCNAME resolves the canonical name for a domain.
func lookupCNAME(domain string) DNSResult {
	start := time.Now()
	cname, err := net.LookupCNAME(domain)
	dur := time.Since(start)

	var results []string
	if err == nil {
		results = []string{cname}
	}

	return DNSResult{
		Domain:     domain,
		RecordType: "CNAME",
		Results:    results,
		Duration:   dur,
		Error:      err,
	}
}

// lookupMX resolves the mail exchange records for a domain.
func lookupMX(domain string) DNSResult {
	start := time.Now()
	mxRecords, err := net.LookupMX(domain)
	dur := time.Since(start)

	var results []string
	if err == nil {
		for _, mx := range mxRecords {
			results = append(results, fmt.Sprintf("%s (priority: %d)", mx.Host, mx.Pref))
		}
	}

	return DNSResult{
		Domain:     domain,
		RecordType: "MX",
		Results:    results,
		Duration:   dur,
		Error:      err,
	}
}

// resolveDomain runs all DNS lookups for a single domain concurrently
// and sends the results to the provided channel.
func resolveDomain(domain string, results chan<- DNSResult, wg *sync.WaitGroup) {
	defer wg.Done()

	var innerWg sync.WaitGroup
	innerWg.Add(3)

	// Launch all three lookup types in parallel
	go func() {
		defer innerWg.Done()
		results <- lookupA(domain)
	}()

	go func() {
		defer innerWg.Done()
		results <- lookupCNAME(domain)
	}()

	go func() {
		defer innerWg.Done()
		results <- lookupMX(domain)
	}()

	innerWg.Wait()
}

// printResult formats and prints a single DNS result with colors.
func printResult(r DNSResult) {
	status := colorGreen + "✓" + colorReset
	if r.Error != nil {
		status = colorRed + "✗" + colorReset
	}

	fmt.Printf("  %s %s%-8s%s │ %s%-25s%s │ ", status, colorYellow, r.RecordType, colorReset, colorCyan, r.Domain, colorReset)

	if r.Error != nil {
		fmt.Printf("%sERROR: %s%s", colorRed, r.Error, colorReset)
	} else if len(r.Results) == 0 {
		fmt.Printf("%sNo records found%s", colorDim, colorReset)
	} else {
		fmt.Printf("%s", strings.Join(r.Results, ", "))
	}

	fmt.Printf(" %s(%s)%s\n", colorDim, r.Duration.Round(time.Microsecond), colorReset)
}

func main() {
	fmt.Printf("\n%s%s╔══════════════════════════════════════════════════════════╗%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║   🌐 Day 50: Go-DNS-Resolver — Concurrent DNS Lookups   ║%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s║              🚀 MILESTONE: 50 DAYS OF GO! 🚀            ║%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s╚══════════════════════════════════════════════════════════╝%s\n\n", colorBold, colorCyan, colorReset)

	// Domains to resolve
	domains := []string{
		"google.com",
		"github.com",
		"golang.org",
		"example.com",
		"cloudflare.com",
	}

	fmt.Printf("%s%s[*] Resolving %d domains concurrently...%s\n\n", colorBold, colorGreen, len(domains), colorReset)

	results := make(chan DNSResult, len(domains)*3) // 3 lookup types per domain
	var wg sync.WaitGroup

	totalStart := time.Now()

	// Launch concurrent resolution for each domain
	for _, domain := range domains {
		wg.Add(1)
		go resolveDomain(domain, results, &wg)
	}

	// Close the results channel once all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results
	var allResults []DNSResult
	for r := range results {
		allResults = append(allResults, r)
	}

	totalDuration := time.Since(totalStart)

	// Print results grouped by domain
	fmt.Printf("  %s%-10s%s │ %s%-25s%s │ %s\n", colorBold, "TYPE", colorReset, colorBold, "DOMAIN", colorReset, colorBold+"RESULTS"+colorReset)
	fmt.Println("  " + strings.Repeat("─", 80))

	currentDomain := ""
	// Group by domain for cleaner output
	for _, domain := range domains {
		for _, r := range allResults {
			if r.Domain != domain {
				continue
			}
			if r.Domain != currentDomain {
				if currentDomain != "" {
					fmt.Println()
				}
				currentDomain = r.Domain
			}
			printResult(r)
		}
	}

	// Summary
	fmt.Println()
	fmt.Println("  " + strings.Repeat("─", 80))
	fmt.Printf("\n  %s%s📊 Summary%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("     Domains resolved : %s%d%s\n", colorGreen, len(domains), colorReset)
	fmt.Printf("     Total lookups    : %s%d%s\n", colorGreen, len(allResults), colorReset)
	fmt.Printf("     Total time       : %s%s%s\n", colorYellow, totalDuration.Round(time.Millisecond), colorReset)

	// Count errors
	errCount := 0
	for _, r := range allResults {
		if r.Error != nil {
			errCount++
		}
	}
	if errCount > 0 {
		fmt.Printf("     Errors           : %s%d%s\n", colorRed, errCount, colorReset)
	} else {
		fmt.Printf("     Errors           : %s0%s\n", colorGreen, colorReset)
	}

	fmt.Printf("\n  %s%s🚀 50 Days of Go — What a journey!%s\n\n", colorBold, colorYellow, colorReset)
}
