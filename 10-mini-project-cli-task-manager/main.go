package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Task represents a single to-do item
type Task struct {
	ID    int    `json:"id"`
	Text  string `json:"text"`
	Done  bool   `json:"done"`
}

// TaskManager handles the list of tasks and persistence
type TaskManager struct {
	Tasks []Task
	mu    sync.Mutex
	file  string
}

// NewTaskManager creates a new manager and loads existing tasks
func NewTaskManager(filename string) *TaskManager {
	tm := &TaskManager{file: filename}
	tm.load()
	return tm
}

func (tm *TaskManager) load() {
	data, err := os.ReadFile(tm.file)
	if err != nil {
		return // File might not exist yet
	}
	json.Unmarshal(data, &tm.Tasks)
}

func (tm *TaskManager) save() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	data, err := json.MarshalIndent(tm.Tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tm.file, data, 0644)
}

func (tm *TaskManager) AddTask(text string) {
	tm.mu.Lock()
	newID := len(tm.Tasks) + 1
	tm.Tasks = append(tm.Tasks, Task{ID: newID, Text: text, Done: false})
	tm.mu.Unlock()
	tm.save()
}

func (tm *TaskManager) ListTasks() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if len(tm.Tasks) == 0 {
		fmt.Println("\n📭 No tasks found!")
		return
	}
	fmt.Println("\n📋 Your Tasks:")
	for _, t := range tm.Tasks {
		status := " "
		if t.Done {
			status = "X"
		}
		fmt.Printf("[%s] %d: %s\n", status, t.ID, t.Text)
	}
}

// backgroundNotifier simulates a background process
func backgroundNotifier(stopChan chan bool) {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Println("\n[BG] Reminder: Keep pushing forward on your Go journey!")
		case <-stopChan:
			return
		}
	}
}

func main() {
	tm := NewTaskManager("tasks.json")
	stopChan := make(chan bool)
	go backgroundNotifier(stopChan)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🚀 Welcome to the Go CLI Task Manager!")
	fmt.Println("---------------------------------------")

	for {
		fmt.Println("\nOptions: [1] Add Task | [2] List Tasks | [3] Exit")
		fmt.Print("Select an option: ")
		
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Print("Enter task description: ")
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			tm.AddTask(text)
			fmt.Println("✅ Task added!")
		case "2":
			tm.ListTasks()
		case "3":
			fmt.Println("👋 Goodbye! See you tomorrow.")
			stopChan <- true
			return
		default:
			fmt.Println("❌ Invalid option, try again.")
		}
	}
}
