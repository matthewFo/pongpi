// +build windows

package pong

import (
	"log"
)

// Type representing a bus connection
type SpiBus struct {
}

func NewSpiBus(busFilePath string, busSpeedHz uint) *SpiBus {
	log.Fatal("Spi not implemented on windows!")
	return nil
}

// Write data to the bus
func (bus *SpiBus) Write(data []byte) (n int, err error) {
	log.Fatal("Spi not implemented on windows!")
	return
}
