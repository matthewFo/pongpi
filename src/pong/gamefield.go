package pong

import (
	"container/list"
)

// Defines all of the information
type GameField struct {

	// Size of the field, from 0 to width exclusive
	width int

	// All of the drawable items, stored in increasing ZIndex order
	drawables *list.List
}

// Initialized a new field
func NewGameField(width int) *GameField {

	return &GameField{
		width:     width,
		drawables: list.New(),
	}
}

// Adds a drawable to the field
func (field *GameField) Add(addDrawable Drawable) {

	curElement := field.drawables.Front()

	if curElement == nil {
		field.drawables.PushFront(addDrawable)
		return
	}

	for ; curElement != nil; curElement = curElement.Next() {

		curDrawable := curElement.Value.(Drawable)

		if addDrawable.ZIndex() < curDrawable.ZIndex() {
			field.drawables.InsertBefore(addDrawable, curElement)
			return
		}
	}

	// got here so it wasn't less than any
	field.drawables.PushBack(addDrawable)
}

// Determines the color at the given position
func (field *GameField) ColorAt(position float64) RGBA {

	color := RGBA{}

	for curElement := field.drawables.Front(); curElement != nil; curElement = curElement.Next() {

		drawable := curElement.Value.(Drawable)

		color = drawable.ColorAt(position, color)
	}

	return color
}

// Animate all Drawables
func (field *GameField) Animate(dt float64) {

	for curElement := field.drawables.Front(); curElement != nil; {

		drawable := curElement.Value.(Drawable)

		if !drawable.Animate(dt) {
			nextElement := curElement.Next()
			field.drawables.Remove(curElement)
			curElement = nextElement
		} else {
			curElement = curElement.Next()
		}
	}
}

// Render each integer position and pass that to the Display
func (field *GameField) RenderTo(display Display) {
	newRender := make([]RGBA, field.width)
	for ledIndex := 0; ledIndex < field.width; ledIndex++ {
		newRender[ledIndex] = field.ColorAt(float64(ledIndex))
	}
	display.Render(newRender)
}

// Returns true if the field of drawables is valid
func (field *GameField) IsValid() bool {

	var prevDrawable Drawable = nil

	for curElement := field.drawables.Front(); curElement != nil; curElement = curElement.Next() {

		curDrawable := curElement.Value.(Drawable)

		if prevDrawable != nil {
			if prevDrawable.ZIndex() > curDrawable.ZIndex() {
				return false
			}
		}

		prevDrawable = curDrawable
	}

	return true
}

// Return number of drawables in the field
func (field *GameField) DrawableLen() int {
	return field.drawables.Len()
}

// Width of the field
func (field *GameField) Width() int {
	return field.width
}
