package draw

import (
	. "pong"
)

// Line that can be drawn on the board
type Line struct {
	// Extends of the line
	leftEdge, rightEdge float64

	color RGBA

	zindex ZIndex
}

var _ Drawable = &Line{}

// Construct a Line
func NewLine(leftEdge, rightEdge float64, color RGBA, zindex ZIndex) *Line {
	return &Line{
		leftEdge:  leftEdge,
		rightEdge: rightEdge,
		color:     color,
		zindex:    zindex,
	}
}

// Returns the color at position blended on top of baseColor
func (line *Line) ColorAt(position float64, baseColor RGBA) RGBA {

	if line.leftEdge <= position && position <= line.rightEdge {
		return line.color.BlendWith(baseColor)
	}

	return baseColor
}

// ZIndex of line
func (line *Line) ZIndex() ZIndex {
	return line.zindex
}

// Animate line
func (line *Line) Animate(dt float64) bool {
	return true
}

// Player that is drawn on the board
type Player struct {
	// line that is drawn
	line *Line

	// if the player is current holding down the button
	paddleActive bool
}

var testPlayer Drawable = &Player{}

// Construct a Line
func NewPlayer(isLeft bool, field *GameField) (player *Player) {

	if isLeft {
		player = &Player{
			line: NewLine(0, 0.75, RGBA{0, 0, 255, 200}, 10),
		}
	} else {
		player = &Player{
			line: NewLine(float64(field.Width())-1.75, float64(field.Width())-1.0, RGBA{0, 255, 0, 200}, 10),
		}
	}

	return
}

// Set if the player is holding down the paddle or not
func (this *Player) UpdatePaddleActive(paddleActive bool) {
	this.paddleActive = paddleActive
}

// Returns the color at position blended on top of baseColor
func (this *Player) ColorAt(position float64, baseColor RGBA) (color RGBA) {

	if !this.paddleActive {
		color = baseColor
	} else {
		color = this.line.ColorAt(position, baseColor)
	}

	return color
}

// ZIndex of the player
func (this *Player) ZIndex() ZIndex {
	return this.line.zindex
}

// Animate player
func (this *Player) Animate(dt float64) bool {
	return true
}
