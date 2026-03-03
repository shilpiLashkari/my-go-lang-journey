# Standard Library Exploration 📚

Go has a "batteries-included" approach, providing a powerful standard library that covers many common needs without requiring external dependencies.

## 📦 Key Packages

### 1. `fmt` (Format)
Used for formatted I/O, similar to C's `printf` and `scanf`.
- `Println`: Prints with a newline.
- `Printf`: Prints with formatting.
- `Sprintf`: Formats a string and returns it.

### 2. `os` (Operating System)
Provides a platform-independent interface to operating system functionality.
- `Args`: Access command-line arguments.
- `Getenv`: Access environment variables.
- `Create`, `Open`, `Write`: File operations.

### 3. `net/http`
Provides HTTP client and server implementations.
- `Get`: Send a GET request.
- `Post`: Send a POST request.
- `ListenAndServe`: Start an HTTP server.

### 4. `encoding/json`
Implements encoding and decoding of JSON.
- `Marshal`: Convert Go data structures to JSON.
- `Unmarshal`: Convert JSON to Go data structures.

### 5. `time`
Provides functionality for measuring and displaying time.
- `Now`: Current time.
- `Sleep`: Pause execution.
- `Format`: Format time as a string.

## 🚀 Try it out!
Check out [main.go](main.go) to see these packages in action!

---

*Ready for more? Check out [Day 10: Building a Mini Project](../10-building-a-mini-project)!*
