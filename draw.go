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
	"errors"
	"fmt"
	"math"
	"os"
)

// RGB represents a color in RGB format with red, green, and blue components
// Each value ranges from 0 to 255
// Used for mapping color names to actual RGB values
// Example: RGB{255, 0, 0} is red
type RGB struct {
	R, G, B int // Values range from 0-255
}

// Color represents a color by its name
// The name must be one of the predefined colors in the ColorMap
// Example: Color{"red"}
type Color struct {
	Name string
}

// Point represents a 2D point in the coordinate system
// x and y are integer coordinates
type Point struct {
	x, y int // x and y coordinates
}

// ColorMap maps color names to RGB values
// The application supports the following colors:
// red, green, blue, yellow, orange, purple, brown, black, white
var ColorMap = map[string]RGB{
	"red":    {255, 0, 0},
	"green":  {0, 255, 0},
	"blue":   {0, 0, 255},
	"yellow": {255, 255, 0},
	"orange": {255, 164, 0},
	"purple": {128, 0, 128},
	"brown":  {165, 42, 42},
	"black":  {0, 0, 0},
	"white":  {255, 255, 255},
}

// Error types defined for different error cases in the application
// errOutOfBounds: Used when a shape or pixel is outside the display
// invalidColor: Used when a color is not in the ColorMap
// fileError: Used when there is a problem creating or writing to a file
var errOutOfBounds = errors.New("Attempt to draw a figure out of bounds of the screen.")
var invalidColor = errors.New("Attempt to use an invalid color.")
var fileError = errors.New("Unable to create PPM file.")

// geometry interface defines methods that all shapes must implement
// draw: Draws the shape on the provided screen
// printShape: Returns a string representation of the shape
type geometry interface {
	// draw draws the shape on the provided screen
	draw(scn screen) (err error)

	// printShape returns a string representation of the shape
	printShape() (s string)
}

// Rectangle struct represents a rectangle defined by lower-left and upper-right points
// ll: Lower-left corner, ur: Upper-right corner, c: Fill color
type Rectangle struct {
	ll Point // Lower-left corner
	ur Point // Upper-right corner
	c  Color // Fill color
}

// Triangle struct represents a triangle defined by three points
// pt0, pt1, pt2: The three vertices, c: Fill color
type Triangle struct {
	pt0 Point // First point
	pt1 Point // Second point
	pt2 Point // Third point
	c   Color // Fill color
}

// Circle struct represents a circle defined by center point and radius
// center: Center point, r: Radius, c: Fill color
type Circle struct {
	center Point // Center point
	r      int   // Radius
	c      Color // Fill color
}

// screen interface defines methods that any display screen must implement
// Used to abstract the display implementation
// initialize: Create a screen with given dimensions
// getMaxXY: Get the maximum x and y dimensions
// drawPixel: Color a pixel at a location
// getPixel: Get the color of a pixel
// clearScreen: Reset all pixels to white
// screenShot: Save the screen to a PPM file
type screen interface {
	initialize(x, y int)
	getMaxXY() (x, y int)
	drawPixel(x, y int, c Color) (err error)
	getPixel(x, y int) (c Color, err error)
	clearScreen()
	screenShot(f string) (err error)
}

// Display struct implements the screen interface
// maxX, maxY: Dimensions of the display
// matrix: 2D slice representing pixel colors
type Display struct {
	maxX   int       // Width of the display
	maxY   int       // Height of the display
	matrix [][]Color // 2D slice representing pixel colors
}

// colorUnknown checks if a color is not defined in the ColorMap
// Returns true if the color is unknown (not in the map)
func colorUnknown(c Color) bool {
	_, exists := ColorMap[c.Name]
	return !exists
}

// outOfBounds checks if a given point would go out of bounds of the screen.
// Returns true if the point is out of bounds, false otherwise.
func outOfBounds(p Point, scn screen) bool {
	xMax, yMax := scn.getMaxXY()
	return p.x < 0 || p.x >= xMax || p.y < 0 || p.y >= yMax
}

// interpolate() is a helper function
// Linearly interpolates between two points (l0, d0) and (l1, d1)
// Returns a slice of integer values representing the interpolated points
func interpolate(l0, d0, l1, d1 int) (values []int) {
	a := float64(d1-d0) / float64(l1-l0)
	d := float64(d0)

	count := l1 - l0 + 1
	for ; count > 0; count-- {
		values = append(values, int(d))
		d = d + a
	}
	return
}

// draw is the Triangle implementation of the geometry.draw method
// Draws a filled triangle using scanline interpolation
// Returns an error if the triangle is out of bounds or if the color is invalid
func (tri Triangle) draw(scn screen) (err error) {
	// Check if drawing this triangle would cause either error
	if outOfBounds(tri.pt0, scn) || outOfBounds(tri.pt1, scn) || outOfBounds(tri.pt2, scn) {
		return errOutOfBounds
	}
	if colorUnknown(tri.c) {
		return invalidColor
	}

	// Sort the points so that y0 <= y1 <= y2
	y0 := tri.pt0.y
	y1 := tri.pt1.y
	y2 := tri.pt2.y
	if y1 < y0 {
		tri.pt1, tri.pt0 = tri.pt0, tri.pt1
	}
	if y2 < y0 {
		tri.pt2, tri.pt0 = tri.pt0, tri.pt2
	}
	if y2 < y1 {
		tri.pt2, tri.pt1 = tri.pt1, tri.pt2
	}
	x0, y0, x1, y1, x2, y2 := tri.pt0.x, tri.pt0.y, tri.pt1.x, tri.pt1.y, tri.pt2.x, tri.pt2.y

	// Interpolate the x-coordinates for the triangle edges
	x01 := interpolate(y0, x0, y1, x1)
	x12 := interpolate(y1, x1, y2, x2)
	x02 := interpolate(y0, x0, y2, x2)

	// Concatenate the short sides
	x012 := append(x01[:len(x01)-1], x12...)

	// Determine which is left and which is right
	var x_left, x_right []int
	m := len(x012) / 2
	if x02[m] < x012[m] {
		x_left = x02
		x_right = x012
	} else {
		x_left = x012
		x_right = x02
	}

	// Draw the horizontal segments (scanlines)
	for y := y0; y <= y2; y++ {
		for x := x_left[y-y0]; x <= x_right[y-y0]; x++ {
			scn.drawPixel(x, y, tri.c)
		}
	}
	return
}

// insideCircle() is a helper function
// Returns true if the tile point is inside the circle with given center and radius
func insideCircle(center, tile Point, r float64) (inside bool) {
	var dx float64 = float64(center.x - tile.x)
	var dy float64 = float64(center.y - tile.y)
	var distance float64 = math.Sqrt(dx*dx + dy*dy)
	return distance <= r
}

// draw is the Rectangle implementation of the geometry.draw method
// It fills in every pixel inside the rectangle with the specified color
// Returns an error if the rectangle is out of bounds or if the color is invalid
func (r Rectangle) draw(scn screen) (err error) {
	// Check if rectangle is out of bounds
	if outOfBounds(r.ll, scn) || outOfBounds(r.ur, scn) {
		return errOutOfBounds
	}
	if colorUnknown(r.c) {
		return invalidColor
	}

	// Fill in rectangle by drawing each pixel (exclusive upper bounds)
	for x := r.ll.x; x < r.ur.x; x++ {
		for y := r.ll.y; y < r.ur.y; y++ {
			err = scn.drawPixel(x, y, r.c)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// draw is the Circle implementation of the geometry.draw method
// Draws a filled circle using the insideCircle helper
// Only draws pixels within the display bounds
// Returns an error if the circle is out of bounds or if the color is invalid
func (c Circle) draw(scn screen) (err error) {
	maxX, maxY := scn.getMaxXY()
	if c.center.x-c.r < 0 || c.center.y-c.r < 0 ||
		c.center.x+c.r >= maxX || c.center.y+c.r >= maxY {
		return errOutOfBounds
	}
	if colorUnknown(c.c) {
		return invalidColor
	}

	// Iterate over the bounding box of the circle
	for y := c.center.y - c.r; y <= c.center.y+c.r; y++ {
		for x := c.center.x - c.r; x <= c.center.x+c.r; x++ {
			if insideCircle(c.center, Point{x, y}, float64(c.r)) {
				if x >= 0 && x < maxX && y >= 0 && y < maxY {
					scn.drawPixel(x, y, c.c)
				}
			}
		}
	}
	return
}

// printShape is the Rectangle implementation of the geometry.printShape method
// Returns a string description of the rectangle with its coordinates
func (r Rectangle) printShape() (s string) {
	return fmt.Sprintf("Rectangle: (%d,%d) to (%d,%d)", r.ll.x, r.ll.y, r.ur.x, r.ur.y)
}

// printShape is the Triangle implementation of the geometry.printShape method
// Returns a string description of the triangle with its coordinates
func (t Triangle) printShape() (s string) {
	return fmt.Sprintf("Triangle: (%d,%d), (%d,%d), (%d,%d)",
		t.pt0.x, t.pt0.y, t.pt1.x, t.pt1.y, t.pt2.x, t.pt2.y)
}

// printShape is the Circle implementation of the geometry.printShape method
// Returns a string description of the circle with its center and radius
func (c Circle) printShape() (s string) {
	return fmt.Sprintf("Circle: centered around (%d,%d) with radius %d",
		c.center.x, c.center.y, c.r)
}

// initialize creates and initializes a display with the specified dimensions
// Sets all pixels to white (the default color)
func (d *Display) initialize(x, y int) {
	d.maxX = x
	d.maxY = y
	d.matrix = make([][]Color, x)
	for i := range d.matrix {
		d.matrix[i] = make([]Color, y)
		for j := range d.matrix[i] {
			d.matrix[i][j] = Color{"white"} // Initialize to white
		}
	}
}

// getMaxXY returns the width and height dimensions of the display
func (d *Display) getMaxXY() (x, y int) {
	return d.maxX, d.maxY
}

// drawPixel sets the color of a pixel at coordinates (x,y)
// Returns errOutOfBounds error if the coordinates are outside the display
// Returns invalidColor error if the specified color is not recognized
func (d *Display) drawPixel(x, y int, c Color) (err error) {
	// Check if pixel is out of bounds
	if x < 0 || y < 0 || x >= d.maxX || y >= d.maxY {
		return errOutOfBounds
	}

	// Check if color is valid
	if colorUnknown(c) {
		return invalidColor
	}

	// Draw the pixel - store directly
	d.matrix[x][y] = c
	return nil
}

// getPixel retrieves the color of a pixel at coordinates (x,y)
// Returns errOutOfBounds error if the coordinates are outside the display
// Returns invalidColor error if the stored color is not recognized
func (d *Display) getPixel(x, y int) (c Color, err error) {
	// Check if pixel is out of bounds
	if x < 0 || y < 0 || x >= d.maxX || y >= d.maxY {
		return Color{}, errOutOfBounds
	}

	// Get the pixel color - retrieve directly
	c = d.matrix[x][y]

	// Check if color is valid
	if colorUnknown(c) {
		return c, invalidColor
	}

	return c, nil
}

// clearScreen resets all pixels in the display to white color
func (d *Display) clearScreen() {
	for i := range d.matrix {
		for j := range d.matrix[i] {
			d.matrix[i][j] = Color{"white"}
		}
	}
}

// screenShot saves the current state of the display to a PPM image file
// The file format follows the P3 PPM format with RGB values
// Returns fileError if there was a problem creating or writing to the file
func (d *Display) screenShot(f string) (err error) {
	file, err := os.Create(f + ".ppm")
	if err != nil {
		return fileError
	}
	defer file.Close()

	// Write header: columns (width) first, then rows (height)
	if _, err = fmt.Fprintf(file, "P3\n%d %d\n255\n", d.maxX, d.maxY); err != nil {
		return fileError
	}

	// Write pixel data row by row, top to bottom
	for y := 0; y < d.maxY; y++ {
		for x := 0; x < d.maxX; x++ {
			color := d.matrix[x][y]
			rgb := ColorMap[color.Name]

			// Write RGB values with space separator, no newline between pixels
			if x > 0 {
				if _, err = fmt.Fprint(file, " "); err != nil {
					return fileError
				}
			}

			if _, err = fmt.Fprintf(file, "%d %d %d", rgb.R, rgb.G, rgb.B); err != nil {
				return fileError
			}
		}

		// Only add newline at the end of each row
		if _, err = fmt.Fprintln(file); err != nil {
			return fileError
		}
	}

	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
