// Package main implements a geometric shape drawing application
// using Go interfaces. This application allows users to draw
// various shapes (rectangles, triangles, circles) of different colors
// on a virtual screen and save the result as a PPM image file.
//
// CS 341, Spring 2025
// Project 5 – Geometry Using Go Interfaces
// Joel Lau Arrieta

package main

import (
	"fmt"
	"strings"
)

// main is the entry point of the program
// It handles user interaction, shape creation, and saving the result to a file
func main() {
	fmt.Println("Project 5: Geometry Using Go Interfaces")
	fmt.Println("CS 341, Spring 2025")
	fmt.Println()
	fmt.Println("This application allows you to draw various shapes")
	fmt.Println("of different colors via interfaces in Go.")
	fmt.Println()

	// Get display dimensions from user
	var rows, cols int
	fmt.Print("Enter the number of rows (x) that you would like the display to have: ")
	fmt.Scan(&rows)
	fmt.Print("Enter the number of columns (y) that you would like the display to have: ")
	fmt.Scan(&cols)
	fmt.Println()

	// Initialize the display
	var d Display
	d.initialize(rows, cols)

	// Drawing loop: repeatedly prompt user to draw shapes until they choose to exit
	for {
		fmt.Println("Select a shape to draw: ")
		fmt.Println("\t R for a rectangle")
		fmt.Println("\t T for a triangle")
		fmt.Println("\t C for a circle")
		fmt.Println(" or X to stop drawing shapes.")

		var choice string
		fmt.Print("Your choice --> ")
		fmt.Scan(&choice)

		// Check if user wants to exit
		if choice == "X" || choice == "x" {
			break
		}

		var shape geometry
		var err error

		// Process user choice and prompt for shape parameters
		switch choice {
		case "R", "r":
			shape, err = drawRectangle()
		case "T", "t":
			shape, err = drawTriangle()
		case "C", "c":
			shape, err = drawCircle()
		default:
			fmt.Println("Invalid choice, please try again.")
			continue
		}

		// Print the shape and attempt to draw it
		fmt.Println(shape.printShape())

		// Draw the shape on the display
		err = shape.draw(&d)
		if err != nil {
			fmt.Printf("**Error: %v\n", err)
		} else {
			shapeName := getShapeName(shape.printShape())
			fmt.Printf("%s drawn successfully.\n", shapeName)
		}
	}

	// Save the drawing to a file
	var filename string
	fmt.Print("Enter the name of the .ppm file in which the results should be saved: ")
	fmt.Scan(&filename)

	err := d.screenShot(filename)
	if err != nil {
		fmt.Printf("**Error: %v\n", err)
	} else {
		fmt.Println("Done. Exiting program...")
	}
}

// getShapeName extracts the shape name from the printShape() output
// Used for user feedback after drawing a shape
func getShapeName(shapeDescription string) string {
	// Find the shape name before the colon
	colonIndex := strings.Index(shapeDescription, ":")
	if colonIndex > 0 {
		return shapeDescription[:colonIndex]
	}
	return "Shape" // Default fallback
}

// drawRectangle prompts the user for rectangle parameters and creates a Rectangle
// Returns a Rectangle object implementing the geometry interface and any error encountered
func drawRectangle() (geometry, error) {
	var llx, lly, urx, ury int
	var colorName string

	fmt.Print("Enter the X and Y values of the lower left corner of the rectangle: ")
	fmt.Scan(&llx, &lly)

	fmt.Print("Enter the X and Y values of the upper right corner of the rectangle: ")
	fmt.Scan(&urx, &ury)

	fmt.Print("Enter the color of the rectangle: ")
	fmt.Scan(&colorName)

	// Create the rectangle
	r := Rectangle{
		ll: Point{llx, lly},
		ur: Point{urx, ury},
		c:  Color{colorName},
	}

	// Check if color is valid
	if colorUnknown(r.c) {
		return r, invalidColor
	}

	return r, nil
}

// drawTriangle prompts the user for triangle parameters and creates a Triangle
// Returns a Triangle object implementing the geometry interface and any error encountered
func drawTriangle() (geometry, error) {
	var x0, y0, x1, y1, x2, y2 int
	var colorName string

	fmt.Print("Enter the X and Y values of the first point of the triangle: ")
	fmt.Scan(&x0, &y0)

	fmt.Print("Enter the X and Y values of the second point of the triangle: ")
	fmt.Scan(&x1, &y1)

	fmt.Print("Enter the X and Y values of the third point of the triangle: ")
	fmt.Scan(&x2, &y2)

	fmt.Print("Enter the color of the triangle: ")
	fmt.Scan(&colorName)

	// Create the triangle
	t := Triangle{
		pt0: Point{x0, y0},
		pt1: Point{x1, y1},
		pt2: Point{x2, y2},
		c:   Color{colorName},
	}

	// Check if color is valid
	if colorUnknown(t.c) {
		return t, invalidColor
	}

	return t, nil
}

// drawCircle prompts the user for circle parameters and creates a Circle
// Returns a Circle object implementing the geometry interface and any error encountered
func drawCircle() (geometry, error) {
	var centerX, centerY, radius int
	var colorName string

	fmt.Print("Enter the X and Y values of the center of the circle: ")
	fmt.Scan(&centerX, &centerY)

	fmt.Print("Enter the value of the radius of the circle: ")
	fmt.Scan(&radius)

	fmt.Print("Enter the color of the circle: ")
	fmt.Scan(&colorName)

	// Create the circle
	c := Circle{
		center: Point{centerX, centerY},
		r:      radius,
		c:      Color{colorName},
	}

	// Check if color is valid
	if colorUnknown(c.c) {
		return c, invalidColor
	}

	return c, nil
}
