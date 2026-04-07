package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// --- ANSI Colors ---
const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorBold   = "\033[1m"
)

type FolderInfo struct {
	Name string
	Size int64
}

func main() {
	targetPath := flag.String("path", ".", "The directory to analyze")
	flag.Parse()

	absPath, _ := filepath.Abs(*targetPath)
	fmt.Printf("%s🔍 Analyzing: %s%s\n\n", colorCyan+colorBold, absPath, colorReset)

	// 1. Read the first level of directories
	entries, err := os.ReadDir(absPath)
	if err != nil {
		fmt.Printf("❌ Error reading directory: %v\n", err)
		return
	}

	var results []FolderInfo
	var totalSize int64

	fmt.Print("⏳ Scanning subdirectories...")
	for _, entry := range entries {
		if entry.IsDir() {
			path := filepath.Join(absPath, entry.Name())
			size := getDirSize(path)
			results = append(results, FolderInfo{Name: entry.Name(), Size: size})
			totalSize += size
		}
	}
	fmt.Print("\r") // Clear the scanning message

	// 2. Sort by size descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Size > results[j].Size
	})

	// 3. Print the report
	if len(results) == 0 {
		fmt.Println("No subdirectories found.")
		return
	}

	maxSize := results[0].Size
	fmt.Printf("%-25s %-10s %s\n", "FOLDER", "SIZE", "VISUAL BREAKDOWN")
	fmt.Println(strings.Repeat("-", 60))

	for _, folder := range results {
		bar := renderBar(folder.Size, maxSize)
		fmt.Printf("%-25s %-10s %s\n", truncate(folder.Name, 24), formatSize(folder.Size), bar)
	}

	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("%-25s %-10s\n", "TOTAL CHECKED", formatSize(totalSize))
}

// getDirSize recursively calculates the size of a directory.
func getDirSize(path string) int64 {
	var size int64 = 0
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip unreadable files
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

// formatSize converts bytes to human-readable strings.
func formatSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// renderBar draws a colored bar based on relative size.
func renderBar(size, maxSize int64) string {
	const barWidth = 20
	if maxSize == 0 {
		return ""
	}
	
	percentage := float64(size) / float64(maxSize)
	filled := int(percentage * float64(barWidth))
	
	color := colorGreen
	if percentage > 0.4 {
		color = colorYellow
	}
	if percentage > 0.8 {
		color = colorRed
	}

	bar := color + strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled) + colorReset
	return fmt.Sprintf("%s %3.0f%%", bar, percentage*100)
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}
