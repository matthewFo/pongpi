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

	expectedColors int
	byteData       []byte
}

var testLedDisplay Display = &LedDisplay{}

// Construct an LedDisplay
func NewLedDisplay(field *GameField, settings SettingsData) *LedDisplay {
	return &LedDisplay{
		bus:            NewSpiBus(settings.SpiFilePath, settings.SpiBusSpeedHz),
		expectedColors: field.Width(),
		byteData:       make([]byte, 4+field.Width()*3+4), // +8 is for the null bytes on front and end of data
	}
}

// Render the colorData to the SPI bus
func (this *LedDisplay) Render(colorData []RGBA) {
	if len(colorData) != this.expectedColors {
		log.Fatal("colorData was not the expected length of ", this.expectedColors, " saw ", len(colorData))
	}

	for colorIndex := 0; colorIndex < len(colorData); colorIndex++ {
		color := colorData[colorIndex]
		byteIndex := colorIndex*3 + 4

		this.byteData[byteIndex+0] = gammaCorrectionLookup[color.G] | 0x80
		this.byteData[byteIndex+1] = gammaCorrectionLookup[color.R] | 0x80
		this.byteData[byteIndex+2] = gammaCorrectionLookup[color.B] | 0x80
	}

	this.bus.Write(this.byteData)
}

// based on a pull request found at http://forums.adafruit.com/viewtopic.php?f=47&t=26591
// is basically precomputing x = pow(i / 255, 3.0) * 127
var gammaCorrectionLookup [256]uint8 = [256]uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 4, 4, 4,
	4, 4, 4, 4, 5, 5, 5, 5, 5, 6, 6, 6, 6, 6, 7, 7,
	7, 7, 7, 8, 8, 8, 8, 9, 9, 9, 9, 10, 10, 10, 10, 11,
	11, 11, 12, 12, 12, 13, 13, 13, 13, 14, 14, 14, 15, 15, 16, 16,
	16, 17, 17, 17, 18, 18, 18, 19, 19, 20, 20, 21, 21, 21, 22, 22,
	23, 23, 24, 24, 24, 25, 25, 26, 26, 27, 27, 28, 28, 29, 29, 30,
	30, 31, 32, 32, 33, 33, 34, 34, 35, 35, 36, 37, 37, 38, 38, 39,
	40, 40, 41, 41, 42, 43, 43, 44, 45, 45, 46, 47, 47, 48, 49, 50,
	50, 51, 52, 52, 53, 54, 55, 55, 56, 57, 58, 58, 59, 60, 61, 62,
	62, 63, 64, 65, 66, 67, 67, 68, 69, 70, 71, 72, 73, 74, 74, 75,
	76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91,
	92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 104, 105, 106, 107, 108,
	109, 110, 111, 113, 114, 115, 116, 117, 118, 120, 121, 122, 123, 125, 126, 127,
}
