package pong

import (
	"log"
	"os"
)

// Type representing a bus connection
type GpioReader struct {
	leftButtonFile  *os.File
	leftPrevious    bool
	rightButtonFile *os.File
	rightPrevious   bool
	data            []byte
}

func NewGpioReader(settings SettingsData) *GpioReader {

	reader := &GpioReader{
		data: make([]byte, 32),
	}

	var err error
	reader.leftButtonFile, err = os.Open(settings.LeftButtonPath)
	if err != nil {
		log.Fatal(err)
	}

	reader.rightButtonFile, err = os.Open(settings.LeftButtonPath)
	if err != nil {
		log.Fatal(err)
	}

	return reader
}

type ButtonEvent int

const (
	NoEvent ButtonEvent = iota
	ButtonPush
	ButtonRelease
)

// get state of the right button
func (this *GpioReader) LeftButton() bool {
	count, err := this.leftButtonFile.Read(this.data)
	if err != nil {
		log.Fatal(err)
	}
	if count != 2 {
		log.Fatal("Expected 2 bytes for left button read and got", count)
	}

	// seek back to beginning of file
	_, err = this.leftButtonFile.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	buttonDown := this.data[0] == 48 // ascii '0'

	return buttonDown

	// if buttonDown != this.leftPrevious {
	// 	this.leftPrevious = buttonDown

	// 	if buttonDown {			
	// 		return ButtonPush
	// 	} else {
	// 		return ButtonRelease
	// 	}
	// }
}

// Get state of the right button
func (this *GpioReader) RightButton() bool {
	count, err := this.rightButtonFile.Read(this.data)
	if err != nil {
		log.Fatal(err)
	}
	if count != 2 {
		log.Fatal("Expected 2 bytes for right button read and got", count)
	}

	// seek back to beginning of file
	_, err = this.rightButtonFile.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	buttonDown := this.data[0] == 48 // ascii '0'

	return buttonDown

	// if buttonDown != this.rightPrevious {
	// 	this.rightPrevious = buttonDown

	// 	if buttonDown {
	// 		return ButtonPush
	// 	} else {
	// 		return ButtonRelease
	// 	}
	// }
}
