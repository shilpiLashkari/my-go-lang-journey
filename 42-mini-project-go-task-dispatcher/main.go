package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job represents a unit of work
type Job struct {
	ID       int
	Payload  string
	Severity int
}

// Result represents the outcome of a job
type Result struct {
	JobID    int
	WorkerID int
	Status   string
	Duration time.Duration
}

// Worker represents a goroutine that processes jobs
type Worker struct {
	ID         int
	JobChannel chan Job
	Quit       chan bool
}

// NewWorker creates a new worker
func NewWorker(id int, jobChan chan Job) *Worker {
	return &Worker{
		ID:         id,
		JobChannel: jobChan,
		Quit:       make(chan bool),
	}
}

// Start initiates the worker's processing loop
func (w *Worker) Start(wg *sync.WaitGroup, results chan<- Result) {
	fmt.Printf("👷 Worker %d: Standing by...\n", w.ID)
	go func() {
		defer wg.Done()
		for {
			select {
			case job := <-w.JobChannel:
				// Simulate processing time based on severity
				processTime := time.Duration(job.Severity*100) * time.Millisecond
				time.Sleep(processTime)

				results <- Result{
					JobID:    job.ID,
					WorkerID: w.ID,
					Status:   "COMPLETED",
					Duration: processTime,
				}
				fmt.Printf("✅ Worker %d: Finished Job %d (Duration: %v)\n", w.ID, job.ID, processTime)

			case <-w.Quit:
				fmt.Printf("💤 Worker %d: Shutting down...\n", w.ID)
				return
			}
		}
	}()
}

// Dispatcher manages the worker pool
type Dispatcher struct {
	WorkerPool chan chan Job
	MaxWorkers int
	JobChannel chan Job
	Workers    []*Worker
}

func NewDispatcher(maxWorkers int, jobChannel chan Job) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		WorkerPool: pool,
		MaxWorkers: maxWorkers,
		JobChannel: jobChannel,
	}
}

func (d *Dispatcher) Run(wg *sync.WaitGroup, results chan<- Result) {
	// 1. Create and start workers
	for i := 1; i <= d.MaxWorkers; i++ {
		worker := NewWorker(i, make(chan Job))
		d.Workers = append(d.Workers, worker)
		wg.Add(1)
		worker.Start(wg, results)
	}

	// 2. Dispatch incoming jobs to available workers
	go func() {
		for {
			select {
			case job := <-d.JobChannel:
				// Simple round-robin dispatch for this demo
				// In a real system, you'd check worker availability via WorkerPool
				worker := d.Workers[job.ID%d.MaxWorkers]
				worker.JobChannel <- job
			}
		}
	}()
}

func main() {
	// Command-line flags
	numWorkers := flag.Int("workers", 5, "Number of workers in the pool")
	numJobs := flag.Int("jobs", 15, "Number of jobs to process")
	queueSize := flag.Int("capacity", 100, "Capacity of the job queue")
	flag.Parse()

	fmt.Printf("🚀 Starting Task Dispatcher (Workers: %d, Jobs: %d, Queue: %d)\n", *numWorkers, *numJobs, *queueSize)
	fmt.Println("---------------------------------------------------------")

	// Channels
	jobChannel := make(chan Job, *queueSize)
	results := make(chan Result, *numJobs)
	var wg sync.WaitGroup

	// Start Dispatcher
	dispatcher := NewDispatcher(*numWorkers, jobChannel)
	dispatcher.Run(&wg, results)

	// Feed jobs into the channel
	go func() {
		for i := 1; i <= *numJobs; i++ {
			job := Job{
				ID:       i,
				Payload:  fmt.Sprintf("Payload Data %d", i),
				Severity: rand.Intn(5) + 1, // 1 to 5
			}
			fmt.Printf("📥 Dispatching Job %d (Severity: %d)...\n", job.ID, job.Severity)
			jobChannel <- job
		}
	}()

	// Collect results
	go func() {
		count := 0
		for range results {
			count++
			if count == *numJobs {
				fmt.Println("---------------------------------------------------------")
				fmt.Printf("🎉 All %d jobs processed successfully!\n", count)
				fmt.Println("Shutting down workers...")
				for _, w := range dispatcher.Workers {
					w.Quit <- true
				}
				return
			}
		}
	}()

	// Keep main alive until all workers finish
	wg.Wait()
	fmt.Println("🛑 System Exit Clean.")
}
