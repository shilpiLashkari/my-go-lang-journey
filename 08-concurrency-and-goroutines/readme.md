# Concurrency in Go ⚡

Concurrency is Go's "killer feature." It allows you to run multiple tasks simultaneously, making your programs faster and more efficient.

## 🏃 Goroutines
A **Goroutine** is a lightweight thread managed by the Go runtime.
- To start one, just add the `go` keyword before a function call.
- They are incredibly cheap (you can run thousands!).

```go
go doSomething()
```

## 🧪 Channels
**Channels** are the pipes that connect goroutines. You can send values into channels from one goroutine and receive those values into another goroutine.
- `ch <- v` // Send `v` to channel `ch`.
- `v := <-ch` // Receive from `ch`, and assign value to `v`.

### WaitGroups
When launching multiple goroutines, I often need to wait for all of them to finish before the program exits. I use `sync.WaitGroup` for this.

## 🚦 The Select Statement
The `select` statement lets a goroutine wait on multiple communication operations. It's like a `switch` but for channels.

```go
select {
case msg1 := <-c1:
    fmt.Println("Received", msg1)
case msg2 := <-c2:
    fmt.Println("Received", msg2)
}
```

## 🚀 Try it out!
Check out [main.go](main.go) to see concurrency in action!

---

*Ready for more? Check out [Day 09: Standard Library Exploration](../09-standard-library-exploration)!*
