package main

import "fmt"

// Define a struct
type Person struct {
	Name string
	Age  int
	City string
}

func main() {
	fmt.Println("--- Pointers ---")
	x := 10
	p := &x // Pointer to x

	fmt.Println("Value of x:", x)
	fmt.Println("Address of x:", p)
	fmt.Println("Value via pointer (*p):", *p)

	*p = 20 // Change value through pointer
	fmt.Println("New value of x after *p = 20:", x)

	fmt.Println("\n--- Structs ---")
	// Creating a struct instance
	p1 := Person{
		Name: "Alice",
		Age:  28,
		City: "New York",
	}
	fmt.Println("Person 1:", p1)
	fmt.Println("Name:", p1.Name)

	// Pointer to a struct
	p2 := &Person{Name: "Bob", Age: 30, City: "London"}
	fmt.Println("Person 2 pointer:", p2)
	fmt.Println("Person 2 Name (auto-dereferenced):", p2.Name)

	fmt.Println("\n--- Structs & Functions (Pass by Value vs Reference) ---")
	
	// Pass by value (copy)
	olderByValue(p1)
	fmt.Println("After olderByValue (copy):", p1.Age)

	// Pass by reference (pointer)
	olderByReference(&p1)
	fmt.Println("After olderByReference (pointer):", p1.Age)
}

func olderByValue(p Person) {
	p.Age += 1
}

func olderByReference(p *Person) {
	p.Age += 1
}
