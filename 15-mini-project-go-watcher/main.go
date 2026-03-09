package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	dir      string
	cmd      string
	exts     string
	interval time.Duration
)

func init() {
	flag.StringVar(&dir, "dir", ".", "Directory to watch")
	flag.StringVar(&cmd, "cmd", "go run main.go", "Command to run on change")
	flag.StringVar(&exts, "ext", ".go,.md", "Comma-separated extensions to watch")
	flag.DurationVar(&interval, "interval", 1*time.Second, "Polling interval")
	flag.Parse()
}

func main() {
	fmt.Printf("🔍 Go-Watcher: Monitoring %s\n", dir)
	fmt.Printf("🚀 Command: %s\n", cmd)
	fmt.Printf("⏱️  Interval: %v\n", interval)

	extensions := strings.Split(exts, ",")
	lastModTimes := make(map[string]time.Time)

	// Initial scan to establish baseline
	lastModTimes = getModTimes(dir, extensions)

	for {
		time.Sleep(interval)
		newModTimes := getModTimes(dir, extensions)

		changed := false
		for path, newTime := range newModTimes {
			if oldTime, ok := lastModTimes[path]; !ok || newTime.After(oldTime) {
				fmt.Printf("\n✨ Change detected in: %s\n", path)
				changed = true
				break
			}
		}

		// Also check for deletions
		if !changed && len(newModTimes) < len(lastModTimes) {
			fmt.Println("\n✨ File deletion detected")
			changed = true
		}

		if changed {
			runCommand(cmd)
			lastModTimes = newModTimes
		}
	}
}

func getModTimes(root string, extensions []string) map[string]time.Time {
	modTimes := make(map[string]time.Time)

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Ignore hidden directories and common noise
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || info.Name() == "vendor") {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			for _, ext := range extensions {
				if strings.HasSuffix(path, ext) {
					modTimes[path] = info.ModTime()
					break
				}
			}
		}
		return nil
	})

	return modTimes
}

func runCommand(commandStr string) {
	fmt.Printf("🔄 Running: %s\n", commandStr)
	fmt.Println(strings.Repeat("-", 20))

	parts := strings.Fields(commandStr)
	if len(parts) == 0 {
		return
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("❌ Command failed: %v\n", err)
	}
	fmt.Println(strings.Repeat("-", 20))
	fmt.Println("✅ Waiting for changes...")
}
