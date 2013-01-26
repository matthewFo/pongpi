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

var testSinusoid Drawable = &Sinusoid{}

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

// ZIndex of line
func (this *Sinusoid) ZIndex() ZIndex {
	return this.zindex
}

// Animate line
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

var testHSLWheel Drawable = &HSLWheel{}

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
