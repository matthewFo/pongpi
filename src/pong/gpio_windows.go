// +build windows

package pong

import ()

// Type representing a bus connection
type GpioReader struct {
}

func NewGpioReader(settings SettingsData) *GpioReader {
	return &GpioReader{}
}

type ButtonEvent int

const (
	NoEvent ButtonEvent = iota
	ButtonPush
	ButtonRelease
)

func (this *GpioReader) LeftButton() bool {
	return false
}

func (this *GpioReader) RightButton() bool {
	return false
}
