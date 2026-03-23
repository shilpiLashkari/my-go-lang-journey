package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// ASCII ramp from darkest to lightest
const asciiRamp = "@#S%?*+;:,. "

func main() {
	imgPath := flag.String("img", "", "Path to the image file (png, jpeg, gif)")
	width := flag.Int("width", 80, "The width of the ASCII output")
	flag.Parse()

	if *imgPath == "" {
		fmt.Println("❌ Please provide an image path using -img")
		os.Exit(1)
	}

	// 1. Open and Decode the image
	file, err := os.Open(*imgPath)
	if err != nil {
		fmt.Printf("❌ Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("❌ Error decoding image: %v\n", err)
		os.Exit(1)
	}

	// 2. Resize the image for terminal display
	// Note: Terminal characters are typically narrow, so we double the width ratio 
	// or halve the height to prevent "tall" ASCII.
	resizedImg := resizeImage(img, *width)

	// 3. Convert pixels to ASCII
	fmt.Println("🎨 Converting to ASCII...\n")
	printASCII(resizedImg)
}

// resizeImage performs a basic Nearest Neighbor downscaling.
func resizeImage(img image.Image, newWidth int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// ASCII characters are taller than they are wide (ratio approx 2:1)
	// We scale height by 0.5 to preserve aspect ratio in the terminal
	newHeight := int(float64(newWidth) * (float64(height) / float64(width)) * 0.5)

	resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Find corresponding pixel in original image
			origX := x * width / newWidth
			origY := y * height / newHeight
			resized.Set(x, y, img.At(origX, origY))
		}
	}

	return resized
}

// printASCII maps pixel brightness to the ASCII ramp.
func printASCII(img image.Image) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			
			// Maps 0-255 (uint8) to 0-len(ramp)-1
			rampIdx := int(c.Y) * (len(asciiRamp) - 1) / 255
			fmt.Print(string(asciiRamp[rampIdx]))
		}
		fmt.Println()
	}
}
