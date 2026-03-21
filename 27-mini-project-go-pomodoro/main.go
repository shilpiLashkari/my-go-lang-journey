package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ANSI colors
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// Session records a completed work or break interval.
type Session struct {
	Type      string
	Duration  time.Duration
	StartedAt time.Time
	Completed bool
}

var sessions []Session

func main() {
	workMins := flag.Int("work", 25, "Work session duration in minutes")
	breakMins := flag.Int("break", 5, "Break session duration in minutes")
	flag.Parse()

	workDur := time.Duration(*workMins) * time.Minute
	breakDur := time.Duration(*breakMins) * time.Minute

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n\n" + colorYellow + "⚠️  Interrupted! Generating session report..." + colorReset)
		cancel()
	}()

	printBanner(*workMins, *breakMins)

	sessionNum := 1
	for {
		// --- Work Session ---
		label := fmt.Sprintf("🍅 Session #%d — Focus Time", sessionNum)
		completed := runTimer(ctx, label, workDur, colorGreen)
		sessions = append(sessions, Session{Type: "Work", Duration: workDur, StartedAt: time.Now(), Completed: completed})
		if !completed {
			break
		}
		fmt.Println(colorGreen + "\n✅ Focus session complete! Time for a break." + colorReset)

		// --- Break Session ---
		label = fmt.Sprintf("☕ Session #%d — Break Time", sessionNum)
		completed = runTimer(ctx, label, breakDur, colorCyan)
		sessions = append(sessions, Session{Type: "Break", Duration: breakDur, StartedAt: time.Now(), Completed: completed})
		if !completed {
			break
		}
		fmt.Println(colorCyan + "\n✅ Break over! Back to work." + colorReset)

		sessionNum++
	}

	printReport()
}

// runTimer shows a live countdown. Returns true if it completed, false if interrupted.
func runTimer(ctx context.Context, label string, duration time.Duration, color string) bool {
	fmt.Printf("\n%s%s%s\n", colorBold, label, colorReset)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	remaining := duration
	for {
		// Print the countdown on the same line using \r
		mins := int(remaining.Minutes())
		secs := int(remaining.Seconds()) % 60
		fmt.Printf("\r  %s⏱  %02d:%02d remaining%s   ", color, mins, secs, colorReset)

		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			remaining -= time.Second
			if remaining <= 0 {
				fmt.Printf("\r  %s⏱  00:00 remaining%s   \n", color, colorReset)
				return true
			}
		}
	}
}

func printBanner(work, brk int) {
	fmt.Println(colorRed + colorBold)
	fmt.Println("  ╔════════════════════════════════╗")
	fmt.Println("  ║     🍅  Go-Pomodoro Timer      ║")
	fmt.Println("  ╚════════════════════════════════╝")
	fmt.Println(colorReset)
	fmt.Printf("  Work: %s%d min%s | Break: %s%d min%s\n", colorGreen, work, colorReset, colorCyan, brk, colorReset)
	fmt.Printf("  Press %sCTRL+C%s at any time to stop.\n\n", colorBold, colorReset)
}

func printReport() {
	fmt.Println("\n" + colorBold + "  📊 SESSION REPORT" + colorReset)
	fmt.Println("  " + "═══════════════════════════════════")

	workCount, breakCount := 0, 0
	for i, s := range sessions {
		status := colorGreen + "✅ Done" + colorReset
		if !s.Completed {
			status = colorYellow + "⚠️  Interrupted" + colorReset
		}
		fmt.Printf("  #%-3d %-8s %s\n", i+1, s.Type, status)
		if s.Type == "Work" && s.Completed {
			workCount++
		} else if s.Type == "Break" && s.Completed {
			breakCount++
		}
	}

	fmt.Println("  " + "═══════════════════════════════════")
	fmt.Printf("  🍅 Completed Focus Sessions: %s%d%s\n", colorGreen, workCount, colorReset)
	fmt.Printf("  ☕ Completed Breaks:         %s%d%s\n", colorCyan, breakCount, colorReset)
	fmt.Println(colorRed + "\n  Great work! 🐹" + colorReset)
}
