# Day 29: Go-Image-To-ASCII - Terminal Art 🖼️✨

## Overview

For Day 29, I built **Go-Image-To-ASCII**, a creative tool that transforms local image files (PNG, JPEG, GIF) into ASCII art directly in the terminal. This project explores pixel-level data manipulation using Go's `image` package. It calculates the brightness (luminance) of every pixel and replaces it with a corresponding character from a custom text ramp.

## Features

- **Multi-Format Support**: Automatically decodes PNG, JPEG, and GIF files.
- **Smart Resizing**: Compresses large images to fit standardized terminal widths without losing too much detail.
- **Brightness Mapping**: Uses an 11-character ramp (`@#S%?*+;:,.`) to represent highlights and shadows.
- **Grayscale Conversion**: Implements standard luminance conversion formulas (`R*0.299 + G*0.587 + B*0.114`).
- **Zero Dependencies**: Purely standard library magic.

## How to Run

```bash
cd 29-mini-project-image-to-ascii

# Convert any local image
go run main.go -img path/to/your/image.jpg

# Set a custom width (default is 80)
go run main.go -img path/to/image.png -width 120
```

## Learning Reflection

- **Image Interpolation**: Implemented a "Nearest Neighbor" resizing algorithm to scale images down for the terminal.
- **Color Luminance**: Learned how the human eye perceives brightness (giving more weight to Green than Blue) and how to calculate grayscale values from RGB.
- **Aspect Ratio Tweak**: Discovered that terminal characters are usually twice as tall as they are wide, so I had to scale the ASCII height by 0.5 to keep the image proportions correct!

---
*Drawing with code, pixel by pixel. 🐹🖼️*
