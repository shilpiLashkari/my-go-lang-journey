package main

import (
	"fmt"
	"math"
)

// Shape interface defines a contract for any type that has an Area method
type Shape interface {
	Area() float64
	Name() string
}

// Rectangle struct
type Rectangle struct {
	Width, Height float64
}

// Area implementation for Rectangle (Method)
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Name() string {
	return "Rectangle"
}

// Circle struct
type Circle struct {
	Radius float64
}

// Area implementation for Circle (Method)
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Name() string {
	return "Circle"
}

// PrintArea is a polymorphic function that accepts any Shape
func PrintArea(s Shape) {
	fmt.Printf("Area of %s: %.2f\n", s.Name(), s.Area())
}

func main() {
	fmt.Println("--- Methods & Interfaces ---")

	r := Rectangle{Width: 10, Height: 5}
	c := Circle{Radius: 7}

	// Using the structs directly
	fmt.Printf("Rectangle Area: %.2f\n", r.Area())
	fmt.Printf("Circle Area: %.2f\n", c.Area())

	fmt.Println("\n--- Polymorphism via Interfaces ---")
	// Using the interface
	shapes := []Shape{r, c}

	for _, shape := range shapes {
		PrintArea(shape)
	}
}
