// +build !linux
package pong

import (
	"log"
)

// Type representing a bus connection
type SpiBus struct {
}

func NewSpiBus(busFilePath string, busSpeedHz uint) *SpiBus {
	log.Fatal("Not implemented on windows!")
}

// Write data to the bus
func (bus *SpiBus) Write(data []byte) (n int, err error) {
	log.Fatal("Not implemented on windows!")
}
