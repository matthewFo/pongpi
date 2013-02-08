package draw

import (
	. "pong"
)

// Player that is drawn on the board
type Player struct {
	start, end              float64
	paddleLeft, paddleRight float64

	lifeColor, paddleColor RGBA

	zindex ZIndex

	// if the player is current holding down the button
	paddleActive bool

	// amount of life left
	life, lifeTotal float64
}

var testPlayer Drawable = &Player{}

// Construct a Line
func NewPlayer(isLeft bool, lifeTime float64, field *GameField) (player *Player) {

	if isLeft {
		player = &Player{
			lifeColor:   RGBA{0, 0, 255, 50},
			paddleColor: RGBA{0, 0, 255, 255},
			zindex:      10,
			start:       0.0,
			end:         (float64(field.Width()) / 2.0) - 1,
			paddleLeft:  0.0,
			paddleRight: 0.5,
			life:        lifeTime,
			lifeTotal:   lifeTime,
		}
	} else {
		player = &Player{
			lifeColor:   RGBA{0, 255, 0, 50},
			paddleColor: RGBA{0, 255, 0, 255},
			zindex:      10,
			start:       float64(field.Width()) - 1.0,
			end:         (float64(field.Width()) / 2.0),
			paddleLeft:  float64(field.Width()) - 1.5,
			paddleRight: float64(field.Width()) - 1.0,
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
	} else if this.paddleActive && left <= position && position <= right {
		color = this.lifeColor.BlendWith(baseColor)
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
func (this *Player) DecreaseLife(dt float64) {
	this.life -= dt
	if this.life < 0 {
		this.life = 0
	}
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
