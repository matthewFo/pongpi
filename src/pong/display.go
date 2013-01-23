package pong

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"time"
)

// A display that can render the field
type Display interface {
	Render([]RGBA)
}

// Web display
type WebDisplay struct {
	previousRender []RGBA
}

var testWebDisplay Display = &WebDisplay{}

// Create a new WebDisplay
func NewWebDisplay(field *GameField) *WebDisplay {
	display := &WebDisplay{
		previousRender: make([]RGBA, int(field.Width())),
	}

	go display.LaunchWebServer()
	return display
}

// Render the field to an internal structure, that can be read out by the webserver
func (this *WebDisplay) Render(data []RGBA) {

	this.previousRender = data
}

// Launches the webserver
func (this *WebDisplay) LaunchWebServer() {

	http.HandleFunc("/", htmlPageHandler)
	http.HandleFunc("/image/", func(w http.ResponseWriter, r *http.Request) { this.imageHandler(w, r) })

	log.Print("Server listening on 8080")
	http.ListenAndServe(":8080", nil)
}

// Serve static html page
func htmlPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
<html>
	<head><script type="text/javascript"><!--
		function reloadpic()
        {
			document.images["gameBoard"].src = "image/test.png";
			setTimeout(reloadpic, 100);
        }
        setTimeout(reloadpic, 100)
	--></script></head>
	<body><img id="gameBoard" src="image/test.png"/></body>
</html>`)

}

// Return newly generating image
func (this *WebDisplay) imageHandler(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-control", "max-age=0, must-revalidate, no-store")

	data := this.previousRender

	spacing := 2
	width, height := len(data), 8
	image := image.NewRGBA(image.Rect(0, 0, width*spacing, height))

	for dataIndex := 0; dataIndex < width; dataIndex++ {

		for y := 0; y < height; y++ {
			image.Set(dataIndex*spacing, y, color.RGBA(data[dataIndex]))
		}
	}

	png.Encode(w, image)
	log.Print("Generated", r.URL, " in", time.Since(startTime))
}

// RGB LED Display
type LedDisplay struct {
	bus *SpiBus
}

var testLedDisplay Display = &LedDisplay{}

// Construct an LedDisplay
func NewLedDisplay(settings SettingsData) *LedDisplay {
	return &LedDisplay{
		bus: NewSpiBus(settings.SpiFilePath, settings.SpiBusSpeedHz),
	}
}

// Render the colorData to the SPI bus
func (this *LedDisplay) Render(colorData []RGBA) {
	byteData := make([]byte, 4+len(colorData)*3+4)

	for colorIndex := 0; colorIndex < len(colorData); colorIndex++ {
		color := colorData[colorIndex]
		byteIndex := colorIndex*3 + 4

		byteData[byteIndex+0] = (color.G >> 1) | 0x80
		byteData[byteIndex+1] = (color.R >> 1) | 0x80
		byteData[byteIndex+2] = (color.B >> 1) | 0x80
	}

	this.bus.Write(byteData)
}
