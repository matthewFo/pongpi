package pong

import (
	"encoding/xml"
	"io/ioutil"
	"log"
)

type SettingsData struct {
	// Max frames per second, app uses thread.sleep to limit FPS
	MaxFPS float64

	// Number of Leds in board
	LedCount int

	// path to the SPI device
	SpiFilePath string

	// speed of the bus
	SpiBusSpeedHz uint

	// Path to the GPIO port for the left button
	LeftButtonPath string

	// Path to the GPIO port for the right button
	RightButtonPath string

	// Min time for a single frame
	MinFrameTime float64 `xml:"-"`
}

// Global settings variable
var Settings SettingsData

// location of the settings file
var settingsFile string = "../settings.xml"

// Read settings from file, setting the global variable
func (settings *SettingsData) Read() {

	fileData, err := ioutil.ReadFile(settingsFile)
	if err != nil {
		log.Fatal(err)
	}
	if err := xml.Unmarshal(fileData, settings); err != nil {
		log.Fatal(err)
	}

	if settings.MaxFPS == 0 {
		settings.MaxFPS = 60
	}

	// setup any derived values
	settings.MinFrameTime = 1.0 / settings.MaxFPS
}

// Write settings to file
func (settings *SettingsData) Write() {
	fileData, err := xml.MarshalIndent(settings, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(settingsFile, fileData, 0777); err != nil {
		log.Fatal(err)
	}
}
