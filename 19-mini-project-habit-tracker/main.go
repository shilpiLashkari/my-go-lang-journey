package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

const habitsFile = "habits.json"

type Habit struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	History   []string  `json:"history"` // Store dates as YYYY-MM-DD
}

type Tracker struct {
	Habits map[string]Habit `json:"habits"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	tracker := loadTracker()
	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: add <habit-name>")
			return
		}
		addHabit(tracker, strings.Join(os.Args[2:], " "))
	case "log":
		if len(os.Args) < 3 {
			fmt.Println("Usage: log <habit-name>")
			return
		}
		logHabit(tracker, strings.Join(os.Args[2:], " "))
	case "list":
		listHabits(tracker)
	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("Usage: remove <habit-name>")
			return
		}
		removeHabit(tracker, strings.Join(os.Args[2:], " "))
	default:
		fmt.Println("Unknown command.")
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: habit-tracker <command> [args]")
	fmt.Println("Commands: add, log, list, remove")
}

func loadTracker() *Tracker {
	file, err := os.ReadFile(habitsFile)
	if err != nil {
		return &Tracker{Habits: make(map[string]Habit)}
	}

	var tracker Tracker
	err = json.Unmarshal(file, &tracker)
	if err != nil {
		return &Tracker{Habits: make(map[string]Habit)}
	}
	return &tracker
}

func saveTracker(tracker *Tracker) {
	data, err := json.MarshalIndent(tracker, "", "  ")
	if err != nil {
		fmt.Printf("❌ Error saving habits: %v\n", err)
		return
	}
	os.WriteFile(habitsFile, data, 0644)
}

func addHabit(tracker *Tracker, name string) {
	if _, ok := tracker.Habits[name]; ok {
		fmt.Printf("⚠️ Habit '%s' already exists.\n", name)
		return
	}

	tracker.Habits[name] = Habit{
		Name:      name,
		CreatedAt: time.Now(),
		History:   []string{},
	}
	saveTracker(tracker)
	fmt.Printf("✅ Added habit: %s\n", name)
}

func logHabit(tracker *Tracker, name string) {
	habit, ok := tracker.Habits[name]
	if !ok {
		fmt.Printf("❓ Habit '%s' not found.\n", name)
		return
	}

	today := time.Now().Format("2006-01-02")
	for _, date := range habit.History {
		if date == today {
			fmt.Printf("ℹ️ Already logged '%s' for today.\n", name)
			return
		}
	}

	habit.History = append(habit.History, today)
	tracker.Habits[name] = habit
	saveTracker(tracker)
	fmt.Printf("🚀 Great job! Logged '%s' for today.\n", name)
}

func listHabits(tracker *Tracker) {
	if len(tracker.Habits) == 0 {
		fmt.Println("📭 No habits tracked yet. Use 'add' to start!")
		return
	}

	fmt.Println("\n🌟 YOUR HABIT PROGRESS")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("%-20s %-8s %s\n", "NAME", "STREAK", "TODAY")

	today := time.Now().Format("2006-01-02")

	for _, habit := range tracker.Habits {
		streak := calculateStreak(habit)
		doneToday := "❌"
		for _, date := range habit.History {
			if date == today {
				doneToday = "✅"
				break
			}
		}
		fmt.Printf("%-20s %-8d %s\n", habit.Name, streak, doneToday)
	}
	fmt.Println(strings.Repeat("-", 40))
}

func calculateStreak(habit Habit) int {
	if len(habit.History) == 0 {
		return 0
	}

	sort.Strings(habit.History)
	streak := 0
	currentTime := time.Now()
	
	// Convert today and yesterday to date-only Time objects
	today, _ := time.Parse("2006-01-02", currentTime.Format("2006-01-02"))
	lastDay := today

	// Work backwards from the most recent entries
	for i := len(habit.History) - 1; i >= 0; i-- {
		entryDate, _ := time.Parse("2006-01-02", habit.History[i])
		
		diff := lastDay.Sub(entryDate).Hours() / 24

		if i == len(habit.History)-1 {
			// If the most recent entry isn't today or yesterday, streak is broken
			if diff > 1 {
				return 0
			}
			streak = 1
		} else {
			if diff == 1 {
				streak++
			} else if diff > 1 {
				break
			}
		}
		lastDay = entryDate
	}

	return streak
}

func removeHabit(tracker *Tracker, name string) {
	if _, ok := tracker.Habits[name]; !ok {
		fmt.Printf("❓ Habit '%s' not found.\n", name)
		return
	}
	delete(tracker.Habits, name)
	saveTracker(tracker)
	fmt.Printf("🗑️ Removed habit: %s\n", name)
}
