package main

import "fmt"

func main() {
	// 1. If/Else - Making Decisions
	fmt.Println("--- 1. If/Else ---")
	age := 18
	if age >= 18 {
		fmt.Println("You are an adult!")
	} else {
		fmt.Println("You are a minor!")
	}

	// 2. Switch - Cleaner Choices
	fmt.Println("\n--- 2. Switch ---")
	day := "Monday"
	switch day {
	case "Monday":
		fmt.Println("Back to work! 💼")
	case "Friday":
		fmt.Println("Weekend is near! 🎉")
	default:
		fmt.Println("Just another day...")
	}

	// 3. For Loops - Doing things over and over
	fmt.Println("\n--- 3. For Loop ---")
	for i := 1; i <= 5; i++ {
		fmt.Println("Count:", i)
	}
}
