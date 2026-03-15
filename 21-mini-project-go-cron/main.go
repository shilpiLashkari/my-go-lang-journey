package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Job represents a single task with an interval.
type Job struct {
	Name     string
	Interval time.Duration
	Task     func()
}

// Scheduler manages the jobs.
type Scheduler struct {
	jobs []Job
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		jobs: []Job{},
	}
}

func (s *Scheduler) Add(job Job) {
	s.jobs = append(s.jobs, job)
	fmt.Printf("✅ Added Job '%s' running every %v\n", job.Name, job.Interval)
}

func (s *Scheduler) Start(ctx context.Context, wg *sync.WaitGroup) {
	fmt.Println("\n🚀 Scheduler is starting all jobs...")
	fmt.Println("----------------------------------------")

	for _, job := range s.jobs {
		wg.Add(1)
		
		// Capture the loop variable
		j := job
		
		go func() {
			defer wg.Done()
			ticker := time.NewTicker(j.Interval)
			defer ticker.Stop()

			// Start immediately on launch
			fmt.Printf("▶️  Executing initial run: %s\n", j.Name)
			j.Task()

			for {
				select {
				case <-ticker.C:
					fmt.Printf("▶️  Executing: %s\n", j.Name)
					j.Task()
				case <-ctx.Done():
					fmt.Printf("⏹️  Stopping Job: %s\n", j.Name)
					return // Exit the goroutine beautifully
				}
			}
		}()
	}
}

func main() {
	fmt.Println("⏰ Go-Cron Prototype")
	fmt.Println("========================================")

	scheduler := NewScheduler()

	// 1. Define Some Sample Jobs
	scheduler.Add(Job{
		Name:     "Fast Health Check",
		Interval: 2 * time.Second,
		Task: func() {
			fmt.Println("    [System is Healthy]")
		},
	})

	scheduler.Add(Job{
		Name:     "Slow Database Backup",
		Interval: 5 * time.Second,
		Task: func() {
			fmt.Println("    [Backing up imaginary database...]")
			time.Sleep(1 * time.Second) // Simulate work
			fmt.Println("    [Backup complete!]")
		},
	})

	// 2. Setup Context for Graceful Shutdown
	// This context will cancel if we receive an OS interrupt signal
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// 3. Catch OS Signals (e.g., CTRL+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// A goroutine watching for the signal
	go func() {
		sig := <-sigChan
		fmt.Printf("\n⚠️  Received signal: %v. Initiating graceful shutdown...\n", sig)
		cancel() // This sends the stop signal to the scheduler's select statements
	}()

	// 4. Start the Engine
	scheduler.Start(ctx, &wg)

	// 5. Wait for Shutdown
	// The main function will block here until wg.Done() is called by all job goroutines
	wg.Wait()

	fmt.Println("========================================")
	fmt.Println("🏁 All jobs stopped cleanly. Goodbye!")
}
