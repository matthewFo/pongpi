package pong

import (
	"image/color"
	"math"
)

// Color used in game
type RGBA color.RGBA

// Defines the z order of different Drawable
type ZIndex int

// Methods required to draw something
type Drawable interface {

	// Computes the color with the given baseColor
	ColorAt(position float64, baseColor RGBA) RGBA

	// The ZIndex of this Drawable thing
	ZIndex() ZIndex

	// Move this Drawable forward in time by dt, returns keepAlive
	Animate(dt float64) (keepAlive bool)
}

// Line that can be drawn on the board
type Line struct {
	// Extends of the line
	leftEdge, rightEdge float64

	color RGBA

	zindex ZIndex
}

var testLine Drawable = Line{}

// Construct a Line
func NewLine(leftEdge, rightEdge float64, color RGBA, zindex ZIndex) Line {
	return Line{
		leftEdge:  leftEdge,
		rightEdge: rightEdge,
		color:     color,
		zindex:    zindex,
	}
}

// Returns the color at position blended on top of baseColor
func (line Line) ColorAt(position float64, baseColor RGBA) RGBA {

	if line.leftEdge < position && position < line.rightEdge {
		return line.color
	}

	return baseColor
}

// ZIndex of line
func (line Line) ZIndex() ZIndex {
	return line.zindex
}

// Animate line
func (line Line) Animate(dt float64) bool {
	return true
}

// Represents
type Sinusoid struct {
	scale float64

	originalColor RGBA

	// offset related to time passing
	offset float64

	zindex ZIndex
}

var testSinusoid Drawable = &Sinusoid{}

// Construct a Sinusoid
func NewSinusoid(scale float64, originalColor RGBA, zindex ZIndex) *Sinusoid {
	return &Sinusoid{
		scale:         scale,
		originalColor: originalColor,
		zindex:        zindex,
	}
}

// Returns the color at position blended on top of baseColor
func (this *Sinusoid) ColorAt(position float64, baseColor RGBA) RGBA {

	sine := (math.Sin((position+this.offset)*this.scale) + 1.0) / 2.0

	return RGBA{
		uint8(float64(this.originalColor.R) * sine),
		uint8(float64(this.originalColor.G) * sine),
		uint8(float64(this.originalColor.B) * sine),
		this.originalColor.A,
	}
}

// ZIndex of line
func (this *Sinusoid) ZIndex() ZIndex {
	return this.zindex
}

// Animate line
func (this *Sinusoid) Animate(dt float64) bool {

	this.offset += dt

	return true
}
