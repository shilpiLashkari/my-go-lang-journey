# Day 19: Go-Habit-Tracker - Consistency Builder 🗓️✨

## Overview

For Day 19, I've built **Go-Habit-Tracker**, a CLI tool designed to help me stay consistent on my journey. After learning about web servers and file systems, I wanted to create something more personal that tracks my progress. This project focuses on handling dates with the `time` package, preserving data with `encoding/json`, and providing a clean, motivational CLI interface.

## Features

- **Personal Habit List**: Easily add or remove habits I want to cultivate.
- **Daily Logging**: Mark habits as completed for the current day.
- **Streak Calculation**: Automatically calculates how many consecutive days I've performed a habit.
- **Visual Progress**: See a quick "Done" status for today across all my habits.
- **Persistence**: Saves my entire habit history to a local `habits.json` file.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 19-mini-project-habit-tracker
   ```
2. Add a new habit (e.g., Coding Go):
   ```bash
   go run main.go add "Coding Go"
   ```
3. Log it for today:
   ```bash
   go run main.go log "Coding Go"
   ```
4. View my current streaks and status:
   ```bash
   go run main.go list
   ```

## Learning Reflection

- **Time Mastery**: Learned how to manipulate dates with `time.Now()` and `time.Parse()`, and how to calculate differences between days to determine streaks.
- **JSON Data Modeling**: Designed a data structure that efficiently stores habit metadata and completion logs.
- **Logical Edge Cases**: Handled logic for what happens when I log a habit multiple times in a day or skip a day.
- **UI Interaction**: Practiced building a user-friendly CLI experience with clear feedback and status indicators.

---
*Stay consistent, stay growing! My Go journey is now more trackable than ever.*
