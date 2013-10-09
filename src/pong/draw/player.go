package draw

import (
	. "pong"
)

// Player that is drawn on the board
type Player struct {

	// start and end of life bar, start is at 0 life, end is full life
	start, end float64

	// bounds of the paddle, used for collision detection
	paddleLeft, paddleRight float64

	// colors that the different parts of the player are drawn
	lifeColor, paddleColor RGBA

	// zindex of player
	zindex ZIndex

	// if the player is current holding down the button
	paddleActive bool

	// amount of life left
	life, lifeTotal float64

	// current amount of animation, goes from 0 to 1 and back
	lifeAnimation float64
}

// rate at which lifeAnimation changes
var lifeAnimationRate float64 = 0.25

var testPlayer Drawable = &Player{}

// Construct a Line
func NewPlayer(isLeft bool, lifeTime float64, field *GameField) (player *Player) {

	if isLeft {
		player = &Player{
			lifeColor:   RGBA{0, 0, 255, 150},
			paddleColor: RGBA{0, 0, 255, 255},
			zindex:      10,
			start:       0.0,
			end:         (float64(field.Width()) / 2.0) - 1,
			paddleLeft:  -0.5,
			paddleRight: 0.5,
			life:        lifeTime,
			lifeTotal:   lifeTime,
		}
	} else {
		player = &Player{
			lifeColor:   RGBA{0, 255, 0, 150},
			paddleColor: RGBA{0, 255, 0, 255},
			zindex:      10,
			start:       float64(field.Width()) - 1.0,
			end:         (float64(field.Width()) / 2.0),
			paddleLeft:  float64(field.Width()) - 1.5,
			paddleRight: float64(field.Width()) - 0.5,
			life:        lifeTime,
			lifeTotal:   lifeTime,
		}
	}

	return
}

// Set if the player is holding down the paddle or not
func (this *Player) UpdatePaddleActive(paddleActive bool) {
	this.paddleActive = paddleActive

	if this.life <= 0.0 {
		this.paddleActive = false
	}
}

// Returns the color at position blended on top of baseColor
func (this *Player) ColorAt(position float64, baseColor RGBA) (color RGBA) {

	lifeBarEnd := (this.life/this.lifeTotal)*(this.end-this.start) + this.start

	left := min(this.start, lifeBarEnd)
	right := max(this.start, lifeBarEnd)

	if this.paddleActive && position == this.start {
		color = this.paddleColor.BlendWith(baseColor)
	} else if left <= position && position <= right && this.life > 0 {

		// animation results in transparency going up and down from 0 to 0.5 when button not pushed, 0.5 to 1 while button pushed
		var alphaAmount = this.lifeAnimation
		if alphaAmount > 0.5 {
			alphaAmount = 1.0 - alphaAmount
		}
		if this.paddleActive {
			alphaAmount += 0.5
		}
		alphaAmount += 0.2
		if alphaAmount > 1.0 {
			alphaAmount = 1.0
		}

		lifeColor := RGBA{this.lifeColor.R, this.lifeColor.G, this.lifeColor.B, uint8(float64(this.lifeColor.A) * alphaAmount)}

		color = lifeColor.BlendWith(baseColor)
	} else {
		color = baseColor
	}

	return
}

// ZIndex of the player
func (this *Player) ZIndex() ZIndex {
	return this.zindex
}

// Animate player
func (this *Player) Animate(dt float64) bool {

	this.lifeAnimation += dt * lifeAnimationRate
	if this.lifeAnimation > 1.0 {
		this.lifeAnimation -= 1.0
	}

	if this.paddleActive {
		this.life -= dt
		if this.life < 0.0 {
			this.life = 0.0
			this.paddleActive = false
		}
	}

	return true
}

// Decrease the amount of life remaining
func (this *Player) DecreaseLife(dt float64) bool {
	this.life -= dt
	if this.life < 0 {
		this.life = 0
		return true
	}

	return false
}

func min(lhs, rhs float64) float64 {
	if lhs < rhs {
		return lhs
	}
	return rhs
}

func max(lhs, rhs float64) float64 {
	if lhs > rhs {
		return lhs
	}
	return rhs
}
