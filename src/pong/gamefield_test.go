package pong

import (
	"testing"
)

// Adding to field should keep things in order
func Test_GameField_Add(t *testing.T) {
	field := NewGameField(100)

	item0 := NewLine(0, 100, Color{255, 255, 255, 255}, 0)
	item1 := NewLine(0, 100, Color{255, 255, 255, 255}, 1)
	item2 := NewLine(0, 100, Color{255, 255, 255, 255}, 2)
	item3 := NewLine(0, 100, Color{255, 255, 255, 255}, 3)
	item4 := NewLine(0, 100, Color{255, 255, 255, 255}, 4)

	field.Add(item0)
	if !field.IsValid() {
		t.Fatal("Invalid field after item0")
	}
	field.Add(item2)
	if !field.IsValid() {
		t.Fatal("Invalid field after item2")
	}
	field.Add(item1)
	if !field.IsValid() {
		t.Fatal("Invalid field after item1")
	}
	field.Add(item3)
	if !field.IsValid() {
		t.Fatal("Invalid field after item3")
	}
	field.Add(item4)
	if !field.IsValid() {
		t.Fatal("Invalid field after item4")
	}
	field.Add(item4)
	if !field.IsValid() {
		t.Fatal("Invalid field after duplicate item4")
	}
}

// Create a Drawable that will kill itself after maxLife amount of time
type CountdownDrawable struct {
	maxLife, curLife float64
}

func (countdown *CountdownDrawable) ColorAt(position float64, baseColor Color) Color {
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
