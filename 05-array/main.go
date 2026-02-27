package main

import "fmt"

func main() {
	fmt.Println("--- Arrays ---")
	// Arrays have a fixed size
	var arr [3]int
	arr[0] = 10
	arr[1] = 20
	arr[2] = 30
	fmt.Println("Array:", arr)
	fmt.Println("Array Length:", len(arr))

	// Array literal
	names := [2]string{"Alice", "Bob"}
	fmt.Println("Names Array:", names)

	fmt.Println("\n--- Slices ---")
	// Slices are dynamic
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("Slice:", nums)
	fmt.Println("Slice Length:", len(nums))
	fmt.Println("Slice Capacity:", cap(nums))

	// Appending to a slice
	nums = append(nums, 6, 7)
	fmt.Println("Slice after append:", nums)

	// Slicing a slice
	subSlice := nums[1:4] // from index 1 to 3
	fmt.Println("Sub Slice [1:4]:", subSlice)

	fmt.Println("\n--- Maps ---")
	// Maps are key-value pairs (like dictionaries or hash maps)
	ages := make(map[string]int)
	ages["Alice"] = 28
	ages["Bob"] = 32
	ages["Charlie"] = 25

	fmt.Println("Map:", ages)
	fmt.Println("Alice's age:", ages["Alice"])

	// Map literal
	scores := map[string]int{
		"Math":    95,
		"English": 88,
		"Science": 92,
	}
	fmt.Println("Scores Map:", scores)

	// Deleting from a map
	delete(scores, "English")
	fmt.Println("Scores Map after deletion:", scores)

	// Checking if a key exists
	score, ok := scores["Math"]
	if ok {
		fmt.Println("Math score found:", score)
	} else {
		fmt.Println("Math score not found")
	}
}
