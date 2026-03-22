# Day 28: Go-Quiz-Master - Timed CSV Quiz Game 🧠⏱️

## Overview

For Day 28, I built **Go-Quiz-Master**, a command-line quiz application. This project focuses on file I/O (reading CSV data), user input management, and concurrency — specifically using a `time.Timer` and channels to handle a global quiz timeout. It's a classic Go exercise that demonstrates how to cleanly handle multiple input sources (stdin and a timer).

## Features

- **CSV Driven**: Loads questions and answers from a simple `problems.csv` file.
- **Timed Sessions**: A global timer (default 30 seconds) starts as soon as you hit Enter.
- **Score Tracking**: Tracks correct and incorrect answers.
- **Concurrent Timeout**: Uses a channel to listen for the timer, ensuring the game ends exactly when time is up, even if you're mid-typing.
- **Customizable**: Set custom time limits and CSV files via flags.

## How to Run

```bash
cd 28-mini-project-quiz-master

# Default: 30 seconds, problems.csv
go run main.go

# Custom: 10 seconds quiz
go run main.go -limit 10
```

## Learning Reflection

- **CSV Parsing**: Used the `encoding/csv` package to read from a file and convert it into a structured slice of problems.
- **Select Statement**: Leveraged `select` to listen to two channels: one for the user's answer and one for the timer's expiration. This is the "Go way" to handle timeouts.
- **Graceful Termination**: Ensures the quiz stops immediately when the timer finishes, printing the final score.

---
*Ready, set, Go! 🐹🧠*
