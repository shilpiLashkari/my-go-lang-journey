# Control Flow in Go 🚦

In this lesson, I'm learning how to control the "flow" of my program using decisions (**If/Else**, **Switch**) and repetitions (**For Loops**).

## 🔀 Making Decisions

### 1. If / Else
The most basic way to make a decision.
```go
if age >= 18 {
    fmt.Println("Adult")
} else {
    fmt.Println("Minor")
}
```

### 2. Switch
A cleaner way to check one variable against many options.
```go
switch day {
case "Monday":
    fmt.Println("Start of the week")
case "Friday":
    fmt.Println("Almost weekend!")
default:
    fmt.Println("Mid-week")
}
```

## 🔄 Doing things again (Loops)

### For Loop
In Go, `for` is the **only** loop! It can be used for counting, iterating over lists, or even acting like a `while` loop.

```go
for i := 1; i <= 5; i++ {
    fmt.Println(i)
}
```

## 🚀 Try it out!
Check out [main.go](main.go) to see these in action!

---

*Ready for more? Check out [Day 04: Functions & Packages](../04-functions-and-packages)!*
