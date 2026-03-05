# Variables and Types in Go 📦

In this lesson, I'm learning how to store information using **Variables**. Think of a variable as a labeled box where you can keep a specific type of item.

## 📦 What are Variables?
A variable is a name for a piece of data that can change. In Go, we tell the computer exactly what "kind" of data (the **Type**) will go into our box.

## 🏷️ Common Data Types
1. **`string`**: For text like `"Hello Go!"`.
2. **`int`**: For whole numbers like `25`.
3. **`bool`**: For true/false values.
4. **`float64`**: For decimal numbers like `1.25`.

## ✍️ How to Create Variables
There are two main ways to create variables in Go:

### 1. The Standard Way (`var`)
Used when you want to be very specific or if you are declaring variables outside of a function.
```go
var name string = "Go Learner"
```

### 2. The Short Way (`:=`)
The most common way inside functions! Go automatically figures out the type for you.
```go
version := 1.25 // Go knows this is a float64
```

## 🚀 Example from `main.go`
In my [main.go](main.go) file, I used these types to print information to the console:

```go
func main() {
    var name string = "Go Learner"
    var age int = 25
    version := 1.25

    fmt.Println("Name:", name)
    fmt.Println("Age:", age)
    fmt.Println("Go Version:", version)
}
```

Happy coding! 🚀
