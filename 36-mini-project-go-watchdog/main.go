package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Watcher struct {
	cmdStr  string
	lastSum map[string]time.Time
	process *exec.Cmd
}

func main() {
	cmd := flag.String("cmd", "", "The terminal command to run on file changes")
	dir := flag.String("dir", ".", "The directory to watch")
	flag.Parse()

	if *cmd == "" {
		fmt.Println("❌ Please provide a command using -cmd")
		fmt.Println("Example: go run main.go -cmd \"echo 'Hello'\"")
		return
	}

	w := &Watcher{
		cmdStr:  *cmd,
		lastSum: make(map[string]time.Time),
	}

	fmt.Printf("🐕 Go-Watchdog: Monitoring %s and running \"%s\"\n", *dir, *cmd)

	// 1. Initial Scan
	w.scan(*dir)
	w.restart()

	// 2. Poll Loop
	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		if changed := w.scan(*dir); changed {
			fmt.Println("\n🔄 Change detected! Restarting...")
			w.restart()
		}
	}
}

// scan recursively checks file modification times.
func (w *Watcher) scan(dir string) bool {
	changed := false
	currentSum := make(map[string]time.Time)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || strings.Contains(path, ".git") {
			return nil
		}

		currentSum[path] = info.ModTime()
		
		// If path is new OR modTime is different
		if oldTime, exists := w.lastSum[path]; !exists || !oldTime.Equal(info.ModTime()) {
			changed = true
		}
		return nil
	})

	// Check if any old keys were deleted
	if len(currentSum) != len(w.lastSum) {
		changed = true
	}

	w.lastSum = currentSum
	return changed
}

// restart kills the previous process and starts a new one.
func (w *Watcher) restart() {
	if w.process != nil && w.process.Process != nil {
		fmt.Printf("🛑 Killing process %d...\n", w.process.Process.Pid)
		w.process.Process.Kill()
		w.process.Wait() // Ensure it's fully cleaned up
	}

	// Prepare the command (Split for multi-part commands like "go run main.go")
	parts := strings.Split(w.cmdStr, " ")
	w.process = exec.Command(parts[0], parts[1:]...)
	w.process.Stdout = os.Stdout
	w.process.Stderr = os.Stderr

	err := w.process.Start()
	if err != nil {
		log.Printf("❌ Failed to start command: %v", err)
	} else {
		fmt.Printf("🚀 Started process %d\n", w.process.Process.Pid)
	}
}
 Lands
