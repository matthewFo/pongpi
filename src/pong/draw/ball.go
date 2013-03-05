package draw

import (
	"math"
	"math/rand"
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

	// the length of the tail of the ball
	tailLength float64

	// if the ball should be hidden this frame or not
	hideBall bool

	// z position of ball
	zindex ZIndex
}

var _ Drawable = &Ball{}

// Construct a Line
func NewBall(field *GameField) *Ball {

	return &Ball{
		position:    0.0,
		velocity:    float64(field.Width()) / 3.0,
		maxPosition: float64(field.Width() - 1),
		tailLength:  8.0,
		zindex:      100,
	}
}

// Returns the color at position blended on top of baseColor
func (this *Ball) ColorAt(position float64, baseColor RGBA) (color RGBA) {

	distance := math.Abs(position - this.position)

	// Add tail flame
	if distance > 0.5 && distance < this.tailLength && ((this.position < position && this.velocity < 0) || (position < this.position && this.velocity > 0)) {

		tailColor := RGBA{255, uint8(rand.Intn(255)), 0, uint8(((this.tailLength - distance) / this.tailLength) * 255.0)}
		baseColor = tailColor.BlendWith(baseColor)
	}

	// Add ball itself as white
	if !this.hideBall && distance < 1 {
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

	return true
}

// Check if the ball went past a player, returns nil or the player that missed the ball
func (this *Ball) MissedByPlayer(leftPlayer, rightPlayer *Player, bounceFactor float64) (missedPlayer *Player) {

	if this.velocity < 0 && this.position < leftPlayer.paddleRight {

		if !leftPlayer.paddleActive && this.position < leftPlayer.paddleLeft {
			// player missed the ball
			PlaySound(MISS)
			return leftPlayer
		} else if leftPlayer.paddleActive {
			// player hit the ball back
			this.position = leftPlayer.paddleRight + (leftPlayer.paddleRight - this.position)
			this.velocity = this.velocity * -bounceFactor
			PlaySound(LEFTBOUNCE)
		}
	} else if this.velocity > 0 && this.position > rightPlayer.paddleLeft {

		if !rightPlayer.paddleActive && this.position > rightPlayer.paddleRight {
			// player missed the ball
			PlaySound(MISS)
			return rightPlayer
		} else if rightPlayer.paddleActive {
			// player hit the ball back
			this.position = rightPlayer.paddleLeft - (this.position - rightPlayer.paddleLeft)
			this.velocity = this.velocity * -bounceFactor
			PlaySound(RIGHTBOUNCE)
		}
	}

	return nil
}

// Reset the position to the middle of the field
func (this *Ball) ResetPosition(field *GameField) {

	startingOffset := 0.25
	if this.velocity < 0 {
		startingOffset = 0.75
	}

	this.position = float64(field.Width()) * startingOffset
}

// Check if the player is doing an offensive hide
func (this *Ball) UpdateOffensiveHide(leftPlayer, rightPlayer *Player) {

	this.hideBall = false
	if this.velocity < 0 && this.position < this.maxPosition/2.0 && rightPlayer.paddleActive {
		this.hideBall = true
	} else if this.velocity > 0 && this.position > this.maxPosition/2.0 && leftPlayer.paddleActive {
		this.hideBall = true
	}
}
