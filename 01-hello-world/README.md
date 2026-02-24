# 👋 Day 01: Hello World

This is the very first step! We're writing a simple "Hello World" program to see how a Go file is put together.

## 📝 The Code (`main.go`)
```go
package main 

import "fmt"

func main() {
    fmt.Println("Hello World!")
}
```

## 🔍 Breaking it Down

*   **`package main`**: Tells Go that this file is an "executable" program (it runs on its own).
*   **`import "fmt"`**: Grabs the `fmt` (Format) tool so we can print text to the screen.
*   **`func main()`**: This is the "Start Button". Everything inside these `{ }` runs first.
*   **`fmt.Println(...)`**: Prints your text and starts a new line.

---

## 🚀 How to Run
1. Open your terminal.
2. Type:
   ```bash
   cd 01-hello-world
   go run main.go
   ```

---

*Ready for the next step? Check out [Day 02: Variables & Types](../02-variable-and-types)!*
