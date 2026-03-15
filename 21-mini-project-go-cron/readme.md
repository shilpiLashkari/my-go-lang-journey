# Day 21: Go-Cron - Lightweight Job Scheduler ⏰

## Overview

For Day 21, I've built **Go-Cron**, a mini background job scheduling system. Real-world backend applications often need to run periodic tasks—like cleaning up old database records, sending daily emails, or polling an API. This project taught me how to use Go's powerful `time` package alongside channels and operating system signals to build a scheduler that runs tasks efficiently and, most importantly, shuts down gracefully without corrupting data.

## Features

- **Interval Execution**: Set functions to run repeatedly at a specified `time.Duration`.
- **Concurrent Jobs**: Multiple tasks run in their own goroutines without blocking each other.
- **Graceful Shutdown**: Catches `CTRL+C` (SIGINT) to allow currently running tasks to finish before exiting the program.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 21-mini-project-go-cron
   ```
2. Start the scheduler:
   ```bash
   go run main.go
   ```
3. Watch the tasks run at different intervals in the console.
4. Press `CTRL+C` to watch the graceful shutdown sequence in action!

## Learning Reflection

- **Contexts for Control**: First time truly understanding how `context.WithCancel` can be passed to multiple goroutines to tell them all to stop at once when the parent context is canceled.
- **Select Statements**: Used the `select` statement to listen simultaneously to a `time.Ticker` channel and a `context.Done()` channel.
- **OS Signals**: Learned how to use `os/signal.Notify` to intercept termination signals from my operating system instead of letting the program crash abruptly.

---
*Scheduling the future, one tick at a time! 🐹⏳*
