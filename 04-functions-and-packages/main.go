package main

import (
	"fmt"
	"math"
)

// 1. Basic function
func sayHello(name string) {
	fmt.Printf("Hello, %s! Welcome to Day 4.\n", name)
}

// 2. Function with a return value
func add(a, b int) int {
	return a + b
}

// 3. Multiple return values
func getCircleInfo(radius float64) (float64, float64) {
	area := math.Pi * radius * radius
	circumference := 2 * math.Pi * radius
	return area, circumference
}

// 4. Named return values
func divide(dividend, divisor float64) (result float64, err string) {
	if divisor == 0 {
		err = "cannot divide by zero"
		return // returns named variables
	}
	result = dividend / divisor
	return
}

func main() {
	fmt.Println("--- Day 4: Functions & Packages ---")

	// Usage 1
	sayHello("Gopher")

	// Usage 2
	sum := add(10, 5)
	fmt.Println("Sum of 10 + 5 =", sum)

	// Usage 3
	a, c := getCircleInfo(5)
	fmt.Printf("Circle with radius 5: Area = %.2f, Circumference = %.2f\n", a, c)

	// Usage 4
	res, e := divide(10, 2)
	fmt.Printf("10 / 2 = %.2f (Error: %s)\n", res, e)

	res2, e2 := divide(10, 0)
	fmt.Printf("10 / 0 = %.2f (Error: %s)\n", res2, e2)
}
