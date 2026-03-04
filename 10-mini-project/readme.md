# Mini Project: CLI Task Manager 📝

Welcome to the final day of the Go journey! To wrap things up, we've built a practical **CLI Task Manager**. This project combines many of the concepts learned over the last 10 days.

## 🚀 Features
- **Add Tasks**: Quickly add new items to your to-do list.
- **List Tasks**: View all your current tasks with their status.
- **Persistence**: Tasks are automatically saved to a `tasks.json` file and loaded when the program starts.
- **Background Monitor**: A simple goroutine that runs in the background to demonstrate concurrency.

## 🧠 Concepts Used
1. **Structs**: To define the `Task` and `TaskManager` data structures.
2. **Standard Library**:
    - `os`: For file handling and environment interactions.
    - `encoding/json`: For persistent data storage.
    - `fmt`: For CLI input/output.
    - `bufio`: For reading user input with spaces.
3. **Concurrency**: A goroutine and channel used for a background status monitor.
4. **Error Handling**: Graceful handling of file I/O and user input errors.

## 🛠️ How to Use
1. Navigate to this directory: `cd 10-mini-project`
2. Run the application: `go run main.go`
3. Follow the on-screen menu to add or view tasks!

---

*Congratulations! You've completed the 10-day Go Lang Journey! 🎈*
