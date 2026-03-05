# Functions & Packages in Go 📦

Go is built around functions. They are the basic building blocks of any Go program. I'm also using packages to organize my code and make it reusable.

## 🛠️ Functions

### 1. Basic Function
A simple function that takes no arguments and returns nothing.
```go
func greet() {
    fmt.Println("Hello Go!")
}
```

### 2. Parameters and Return Types
Functions can take input and return a value.
```go
func add(a int, b int) int {
    return a + b
}
```

### 3. Multiple Return Values
One of Go's coolest features! A function can return more than one thing.
```go
func swap(a, b string) (string, string) {
    return b, a
}
```

### 4. Named Return Values
You can name the return values in the function signature.
```go
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return // "Naked" return
}
```

## 📦 Packages

Packages are Go's way of organizing code.
- Every Go file starts with `package <name>`.
- The `main` package is special; it tells Go that this is an executable program.
- Use `import` to bring in code from other packages.

```go
import "fmt"
import "math"
```

## 🚀 Try it out!
Check out [main.go](main.go) to see functions in action!

---

*Ready for more? Check out [Day 05: Data Structures](../05-data-structures)!*
