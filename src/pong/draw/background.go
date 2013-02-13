package draw

import (
	"math"
	. "pong"
)

// Represents a background animation of a sinusoid moving forward
type Sinusoid struct {

	// length of field
	scale float64

	// offset related to time passing, from 0 to 1
	offsets [3]float64

	zindex ZIndex

	sineLookup []uint8
}

var _ Drawable = &Sinusoid{}

// Construct a Sinusoid
func NewSinusoid(field *GameField, zindex ZIndex) *Sinusoid {
	sine := &Sinusoid{
		scale:   float64(field.Width()),
		offsets: [3]float64{0.0, 0.0, 0.0},
		zindex:  zindex,
	}

	sine.buildLookup()

	return sine
}

// build lookup table to make rendering much faster by precomputing math.Sin
func (this *Sinusoid) buildLookup() {
	this.sineLookup = make([]uint8, 256)
	for index := 0; index < 256; index++ {
		fieldPercentage := float64(index) / 256

		value := (math.Sin(fieldPercentage*2.0*math.Pi+this.offsets[0]) + 1.0) / 2.0
		this.sineLookup[index] = uint8(value*255.0) >> 1
	}
}

// lookup the sine value instead of computing using math.Sin
func (this *Sinusoid) lookup(fieldPercentage float64) uint8 {
	if fieldPercentage > 1 {
		fieldPercentage -= 1
	}
	return this.sineLookup[int(fieldPercentage*256)]
}

// Returns the color at position blended on top of baseColor
func (this *Sinusoid) ColorAt(position float64, baseColor RGBA) RGBA {

	// 0 to 1
	fieldPercentage := position / this.scale

	return RGBA{
		this.lookup(fieldPercentage + this.offsets[0]),
		this.lookup(fieldPercentage + this.offsets[1]),
		this.lookup(fieldPercentage + this.offsets[2]),
		255,
	}
}

// ZIndex
func (this *Sinusoid) ZIndex() ZIndex {
	return this.zindex
}

// Animate
func (this *Sinusoid) Animate(dt float64) bool {

	this.offsets[0] += dt * 0.27
	if this.offsets[0] > 1 {
		this.offsets[0] -= 1
	}

	this.offsets[1] += dt * 0.41
	if this.offsets[1] > 1 {
		this.offsets[1] -= 1
	}

	this.offsets[2] += dt * 0.59
	if this.offsets[2] > 1 {
		this.offsets[2] -= 1
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

var _ Drawable = &HSLWheel{}

// Construct an HSLWheel
func NewHSLWheel(field *GameField, zindex ZIndex) *HSLWheel {
	return &HSLWheel{
		scale:  float64(field.Width()),
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

// ZIndex
func (this *HSLWheel) ZIndex() ZIndex {
	return this.zindex
}

// Animate
func (this *HSLWheel) Animate(dt float64) bool {

	this.hue += dt * 0.1

	if this.hue > 1.0 {
		this.hue -= 1.0
	}

	//log.Print("Hue", this.hue)

	return true
}

// Represents a background animation of moving through the HSL color space
type StepFunction struct {

	// width of step function
	center float64

	// scale of the steps
	scale float64

	// 255.0 / scale
	scaleInverse float64

	// time offset
	offset float64

	// color at max alpha
	baseColor RGBA

	zindex ZIndex
}

var _ Drawable = &StepFunction{}

// Construct a new StepFunction
func NewStepFunction(center, stepSize float64, baseColor RGBA, zindex ZIndex) *StepFunction {
	return &StepFunction{
		center:       center,
		scale:        stepSize,
		scaleInverse: 255.0 / stepSize,
		offset:       0.0,
		baseColor:    baseColor,
		zindex:       zindex,
	}
}

// Returns the color at position blended on top of baseColor
func (this *StepFunction) ColorAt(position float64, baseColor RGBA) RGBA {

	position = math.Abs(position - this.center)

	alpha := uint8(math.Mod(position+this.offset, this.scale) * this.scaleInverse)

	color := RGBA{
		this.baseColor.R,
		this.baseColor.G,
		this.baseColor.B,
		alpha,
	}

	return color.BlendWith(baseColor)
}

// ZIndex
func (this *StepFunction) ZIndex() ZIndex {
	return this.zindex
}

// Animate
func (this *StepFunction) Animate(dt float64) bool {

	this.offset += dt * 40.0

	if this.offset > this.scale {
		this.offset -= this.scale
	}

	return true
}

// Represents a background animation of moving through the HSL color space
type Countdown struct {

	// total time counted so far
	time float64

	// length of entire countdown
	totalTime float64
}

var _ Drawable = &Countdown{}

// Construct a new StepFunction
func NewCountdown(field *GameField, totalTime float64) *Countdown {
	return &Countdown{
		time:      0.0,
		totalTime: totalTime,
	}
}

// Returns the color at position blended on top of baseColor
func (this *Countdown) ColorAt(position float64, baseColor RGBA) RGBA {

	greenFactor := this.time / this.totalTime
	redFactor := 1.0 - greenFactor

	if greenFactor > 0.8 {
		return baseColor
	}

	return RGBA{
		uint8(redFactor * 255),
		uint8(greenFactor * 255),
		0,
		255,
	}
}

// ZIndex
func (this *Countdown) ZIndex() ZIndex {
	return 0
}

// Animate
func (this *Countdown) Animate(dt float64) bool {

	this.time += dt

	if this.time >= this.totalTime {
		this.time = this.totalTime
	}

	return true
}

// Amount of time remaining in countdown
func (this *Countdown) TimeRemaining() float64 {
	return this.totalTime - this.time
}

// Represents a background animation of moving through the HSL color space
type Winner struct {

	// Bounds of color to draw
	left, right float64

	// color to be flashed on / off
	color RGBA

	// total time counted so far
	time float64

	// length of entire countdown
	totalTime float64
}

var _ Drawable = &Winner{}

// Construct a new StepFunction
func NewWinner(field *GameField, leftWon bool, totalTime float64) *Winner {

	if !leftWon {
		return &Winner{
			time:      0.0,
			totalTime: totalTime,
			left:      0,
			right:     (float64(field.Width()) / 2.0) - 1,
			color:     RGBA{0, 0, 255, 255},
		}
	} else {
		return &Winner{
			time:      0.0,
			totalTime: totalTime,
			left:      (float64(field.Width()) / 2.0),
			right:     float64(field.Width()) - 1,
			color:     RGBA{0, 255, 0, 255},
		}
	}

	return nil
}

// Returns the color at position blended on top of baseColor
func (this *Winner) ColorAt(position float64, baseColor RGBA) RGBA {

	if this.left <= position && position <= this.right && int(this.time*4)%2 == 0 {
		return this.color
	}

	return baseColor
}

// ZIndex
func (this *Winner) ZIndex() ZIndex {
	return 0
}

// Animate
func (this *Winner) Animate(dt float64) bool {

	this.time += dt

	if this.time >= this.totalTime {
		this.time = this.totalTime
	}

	return true
}

// Amount of time remaining in countdown
func (this *Winner) TimeRemaining() float64 {
	return this.totalTime - this.time
}
