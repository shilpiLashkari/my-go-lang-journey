# 👋 Day 01: Hello World

This is the very first step of my Go journey! Here, I've written a simple "Hello World" program to understand the basic structure of a Go file.

## 📝 The Code (`main.go`)

```go
package main 

import "fmt"

func main(){
	fmt.Println("Hello World!");
}
```

## 🔍 BreakDown

### 1. `package main`
Every Go file must start with a package declaration. The `main` package is a special package that tells the Go compiler that this file should compile as an executable program rather than a shared library.

### 2. `import "fmt"`
The `import` keyword is used to include code from other packages. Here, I'm importing the `fmt` (format) package, which contains functions for formatting text, including printing to the console.

### 3. `func main()`
The `main` function is the entry point of the program. When you run the execution, the code inside the curly braces `{}` of the `main` function is what gets executed first.

### 4. `fmt.Println(...)`
This is a function from the `fmt` package. `Println` stands for "Print Line". It outputs the text followed by a new line.

---

## 🚀 How to Run

1. Open your terminal.
2. Navigate to this directory:
   ```bash
   cd 01-hello-world
   ```
3. Run the program:
   ```bash
   go run main.go
   ```

---

*Onward to [Day 02: Variables & Types](../README.md#-roadmap)!*
