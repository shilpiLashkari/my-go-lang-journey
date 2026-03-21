# Day 27: Go-Pomodoro - Terminal Focus Timer 🍅⏱️

## Overview

For Day 27, I built **Go-Pomodoro**, a terminal-based Pomodoro timer. The Pomodoro Technique is a productivity method where you work for 25 minutes, then take a 5-minute break. This project brings it to the terminal with a live countdown, colored output, session history tracking, and graceful CTRL+C handling.

## Features

- **Live Countdown**: Ticks every second and re-renders inline using `\r` (carriage return) — no terminal library needed.
- **Work & Break Sessions**: Automatically alternates between 25-min focus and 5-min break intervals.
- **Session History**: Logs completed sessions with timestamps.
- **Configurable Durations**: Use `-work` and `-break` flags to customize timings.
- **Graceful Exit**: Press CTRL+C at any time to see your full session log before quitting.

## How to Run

```bash
cd 27-mini-project-go-pomodoro

# Default: 25 min work, 5 min break
go run main.go

# Custom: 2 min work, 1 min break (for testing)
go run main.go -work 2 -break 1
```

## Learning Reflection

- **Inline Terminal Updates**: Used `\r` (carriage return without newline) and `fmt.Printf` to overwrite the same line every second — a simple animation technique.
- **`time.Ticker` vs `time.Sleep`**: Chose `time.Ticker` for accurate per-second ticks that don't drift from sleep delays.
- **OS Signal + Context**: Combined `os/signal` with a cancellable `context` to cleanly interrupt any countdown and print history.

---
*25 minutes of Go, 5 minutes of tea. 🐹🍅*
