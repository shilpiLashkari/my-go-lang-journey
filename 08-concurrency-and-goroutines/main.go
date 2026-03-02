package main

import (
	"fmt"
	"sync"
	"time"
)

// A simple worker function
func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // Notify WaitGroup that this worker is done

	fmt.Printf("Worker %d starting...\n", id)
	time.Sleep(time.Second) // Simulate work
	fmt.Printf("Worker %d done!\n", id)
}

func main() {
	fmt.Println("--- Day 8: Concurrency & Goroutines ---")

	// 1. Basic Goroutine & WaitGroup
	var wg sync.WaitGroup
	fmt.Println("\nLaunching workers...")
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	wg.Wait() // Wait for all workers to finish
	fmt.Println("All workers finished!")

	// 2. Channels - Communication between Goroutines
	fmt.Println("\n--- Using Channels ---")
	messageChannel := make(chan string)

	go func() {
		time.Sleep(500 * time.Millisecond)
		messageChannel <- "Hello from the channel! 🚀"
	}()

	msg := <-messageChannel
	fmt.Println("Received message:", msg)

	// 3. Select Statement - Multiplexing Channels
	fmt.Println("\n--- Using Select ---")
	chan1 := make(chan string)
	chan2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		chan1 <- "Message from Channel 1"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		chan2 <- "Message from Channel 2"
	}()

	for i := 0; i < 2; i++ {
		select {
		case m1 := <-chan1:
			fmt.Println(m1)
		case m2 := <-chan2:
			fmt.Println(m2)
		}
	}
	fmt.Println("\nDay 8 complete! Concurrency is powerful.")
}
