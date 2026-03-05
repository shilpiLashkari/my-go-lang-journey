package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

// ProjectStats holds the data I'm collecting
type ProjectStats struct {
	GoFiles    int32
	MDFiles    int32
	TotalLines int64
}

func main() {
	fmt.Println("📊 Go Journey Analytics Tool")
	fmt.Println("---------------------------")

	// I'll start from the parent directory to scan the whole repository
	root := ".."
	stats := ProjectStats{}
	var wg sync.WaitGroup

	// Scan the repository
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// I'll ignore the .git folder and the current project folder to avoid recursion
		if info.IsDir() && (info.Name() == ".git" || info.Name() == "11-mini-project-journey-analytics") {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			ext := filepath.Ext(path)
			if ext == ".go" {
				atomic.AddInt32(&stats.GoFiles, 1)
				wg.Add(1)
				go countLinesConcurrent(path, &stats, &wg)
			} else if strings.ToLower(ext) == ".md" {
				atomic.AddInt32(&stats.MDFiles, 1)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Error scanning repository: %v\n", err)
		return
	}

	wg.Wait()

	// Display the dashboard
	fmt.Println("\n🚀 MY LEARNING DASHBOARD")
	fmt.Println("========================")
	fmt.Printf("📂 Total Go Files:    %d\n", stats.GoFiles)
	fmt.Printf("📝 Total Readme Files: %d\n", stats.MDFiles)
	fmt.Printf("💻 Total Lines of Go:  %d\n", stats.TotalLines)
	fmt.Println("========================")
	fmt.Println("Keep it up! Each line is a step forward. 🐹")
}

func countLinesConcurrent(path string, stats *ProjectStats, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines int64
	for scanner.Scan() {
		lines++
	}

	atomic.AddInt64(&stats.TotalLines, lines)
}
