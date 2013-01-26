package draw

import (
	"math"
	. "pong"
)

// Player that is drawn on the board
type Ball struct {

	// current position of the ball
	position float64

	// direction and speed of the ball in leds / second
	velocity float64

	// max position of ball, min is 0
	maxPosition float64

	// z position of ball
	zindex ZIndex
}

var testBall Drawable = &Ball{}

// Construct a Line
func NewBall(field *GameField) *Ball {

	return &Ball{
		position:    0,
		velocity:    float64(field.Width()) / 3.0,
		maxPosition: float64(field.Width() - 1),
		zindex:      100,
	}
}

// Returns the color at position blended on top of baseColor
func (this *Ball) ColorAt(position float64, baseColor RGBA) (color RGBA) {

	distance := math.Abs(position - this.position)
	if distance < 1 {
		color = RGBA{255, 255, 255, uint8((1.0 - distance) * 255.0)}
		color = color.BlendWith(baseColor)
	} else {
		color = baseColor
	}

	return color
}

// ZIndex of the ball
func (this *Ball) ZIndex() ZIndex {
	return this.zindex
}

// Animate ball
func (this *Ball) Animate(dt float64) bool {
	this.position += this.velocity * dt
	if this.position > this.maxPosition {
		this.position = this.maxPosition - (this.position - this.maxPosition)
		this.velocity = -this.velocity
	} else if this.position < 0 {
		this.position = -this.position
		this.velocity = -this.velocity
	}

	//log.Print(this.position, this.velocity)

	return true
}
