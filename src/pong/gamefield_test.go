package pong

import (
	"testing"
)

// Create a Drawable that will kill itself after maxLife amount of time
type CountdownDrawable struct {
	maxLife, curLife float64
}

func (countdown *CountdownDrawable) ColorAt(position float64, baseColor RGBA) RGBA {
	return baseColor
}

func (countdown *CountdownDrawable) ZIndex() ZIndex {
	return 1
}

func (countdown *CountdownDrawable) Animate(dt float64) (keepAlive bool) {
	countdown.curLife += dt
	return countdown.curLife < countdown.maxLife
}

// Animate call to a Drawable that returns false should result in that drawable being removed from the field
func Test_GameField_AddAnimate(t *testing.T) {
	field := NewGameField(100)

	items := []Drawable{
		&CountdownDrawable{maxLife: 1},
		&CountdownDrawable{maxLife: 2},
		&CountdownDrawable{maxLife: 4},
		&CountdownDrawable{maxLife: 3},
		&CountdownDrawable{maxLife: 5},
		&CountdownDrawable{maxLife: 6},
	}

	for _, item := range items {
		field.Add(item)
	}

	if !field.IsValid() {
		t.Fatal("Invalid field")
	}
	Assert(field.DrawableLen(), len(items), "Initial length", t)

	for iter := 1; iter <= len(items); iter++ {

		field.Animate(1.0)

		if !field.IsValid() {
			t.Fatal("Invalid field")
		}

		Assert(field.DrawableLen(), len(items)-iter, "Length after Animate", t)
	}

	Assert(field.DrawableLen(), 0, "Field should be empty", t)
}

// Helper assert method
func Assert(actual, expected int, message string, t *testing.T) {
	if actual != expected {
		t.Fatal(message, actual, "vs expected", expected)
	}
}
