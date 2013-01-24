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

func blend(foreground, background RGBA) (color RGBA) {

	fr,fg,fb,fa := uint(foreground.R), uint(foreground.G), uint(foreground.B), uint(foreground.A)
	br,bg,bb,ba := uint(background.R), uint(background.G), uint(background.B), uint(background.A)

	opacity := fa
	backgroundOpacity := (ba * (255 - fa)) >> 8

	newColor := RGBA{
		uint8((fr * opacity) >> 8 + (br * backgroundOpacity) >> 8),
		uint8((fg * opacity) >> 8 + (bg * backgroundOpacity) >> 8),
		uint8((fb * opacity) >> 8 + (bb * backgroundOpacity) >> 8),
		uint8(opacity),
	}

	return newColor
}

// Line that can be drawn on the board
type Line struct {
	// Extends of the line
	leftEdge, rightEdge float64

	color RGBA

	zindex ZIndex
}

var testLine Drawable = &Line{}

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
		return blend(line.color, baseColor)
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
	visible bool
}

var testPlayer Drawable = &Player{}

// Construct a Line
func NewPlayer(isLeft bool, field *GameField) (player *Player) {

	if isLeft {
		player = &Player{
			line: NewLine(0, float64(field.Width()/2), RGBA{255, 0, 0, 200}, 10),
		}
	} else {
		player = &Player{
			line: NewLine(float64(field.Width()/2)+1, float64(field.Width()), RGBA{0, 255, 0, 200}, 10),
		}
	}

	return
}

// Set if the player is visible or not
func (this *Player) UpdateVisible(visible bool) {
	this.visible = visible
}

// Returns the color at position blended on top of baseColor
func (this *Player) ColorAt(position float64, baseColor RGBA) (color RGBA) {

	if !this.visible {
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

// Represents a background animation of a sinusoid moving forward
type Sinusoid struct {

	// length of field
	scale float64

	// offset related to time passing
	offsets [3]float64

	zindex ZIndex
}

var testSinusoid Drawable = &Sinusoid{}

// Construct a Sinusoid
func NewSinusoid(field *GameField, zindex ZIndex) *Sinusoid {
	return &Sinusoid{
		scale:   float64(field.Width()),
		offsets: [3]float64{0.0, 0.0, 0.0},
		zindex:  zindex,
	}
}

// Returns the color at position blended on top of baseColor
func (this *Sinusoid) ColorAt(position float64, baseColor RGBA) RGBA {

	// 0 to 1
	fieldPercentage := position / this.scale

	red := (math.Sin(fieldPercentage*2.0*math.Pi+this.offsets[0]) + 1.0) / 2.0
	green := (math.Sin(fieldPercentage*2.0*math.Pi+this.offsets[1]) + 1.0) / 2.0
	blue := (math.Sin(fieldPercentage*2.0*math.Pi+this.offsets[2]) + 1.0) / 2.0

	//sine := (math.Sin((position+this.offset)*this.scale) + 1.0) / 2.0

	return RGBA{
		uint8(red * 255.0),
		uint8(green * 255.0),
		uint8(blue * 255.0),
		255,
	}
}

// ZIndex of line
func (this *Sinusoid) ZIndex() ZIndex {
	return this.zindex
}

// Animate line
func (this *Sinusoid) Animate(dt float64) bool {

	this.offsets[0] += dt * 1.62
	if this.offsets[0] > 2.0*math.Pi {
		this.offsets[0] -= 2.0 * math.Pi
	}

	this.offsets[1] += dt * 2.49
	if this.offsets[1] > 2.0*math.Pi {
		this.offsets[1] -= 2.0 * math.Pi
	}

	this.offsets[2] += dt * 3.58
	if this.offsets[2] > 2.0*math.Pi {
		this.offsets[2] -= 2.0 * math.Pi
	}

	return true
}

// Represents a background animation of moving through the HSL color space
type HSLWheel struct {

	// Hwo to scale lumniosity so that it's event spread across all points
	scale float64

	// goes from 0 to 1 and then wraps
	hue float64

	zindex ZIndex
}

var testHSLWheel Drawable = &HSLWheel{}

// Construct an HSLWheel
func NewHSLWheel(field *GameField, zindex ZIndex) *HSLWheel {
	return &HSLWheel{
		scale:  float64(field.width),
		hue:    0.0,
		zindex: zindex,
	}
}

// Returns the color at position blended on top of baseColor
func (this *HSLWheel) ColorAt(position float64, baseColor RGBA) RGBA {

	luminosity := position / this.scale

	// shift it up because we don't care much about the very dark colors
	//luminosity = luminosity*0.8 + 0.2

	return hslToRGB(this.hue, 1.0, luminosity)
}

// Convert HSL to RGB, based on http://mjijackson.com/2008/02/rgb-to-hsl-and-rgb-to-hsv-color-model-conversion-algorithms-in-javascript
func hslToRGB(hue, saturation, luminosity float64) RGBA {

	var red, green, blue uint8

	if saturation == 0 {
		red = 255
		green = 255
		blue = 255
	} else {

		hueToRGB := func(p, q, t float64) float64 {
			if t < 0.0 {
				t += 1.0
			}
			if t > 1.0 {
				t -= 1.0
			}

			if t < 1.0/6.0 {
				return p + (q-p)*6.0*t
			}
			if t < 1.0/2.0 {
				return q
			}
			if t < 2.0/3.0 {
				return p + (q-p)*(2.0/3.0-t)*6.0
			}
			return p
		}

		var q float64
		if luminosity < 0.5 {
			q = luminosity * (1.0 + saturation)
		} else {
			q = luminosity + saturation - luminosity*saturation
		}

		p := 2*luminosity - q
		red = uint8(hueToRGB(p, q, hue+1.0/3.0) * 255.0)
		green = uint8(hueToRGB(p, q, hue) * 255.0)
		blue = uint8(hueToRGB(p, q, hue-1.0/3.0) * 255.0)
	}

	//log.Print("Converted ", hue, saturation, luminosity, " to ", red, green, blue)

	return RGBA{red, green, blue, 255}
}

// ZIndex of line
func (this *HSLWheel) ZIndex() ZIndex {
	return this.zindex
}

// Animate line
func (this *HSLWheel) Animate(dt float64) bool {

	this.hue += dt * 0.1

	if this.hue > 1.0 {
		this.hue -= 1.0
	}

	//log.Print("Hue", this.hue)

	return true
}
