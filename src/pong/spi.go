// +build !windows

package pong

import (
	"log"
	"os"
	"syscall"
	"unsafe"
)

// Type representing a bus connection
type SpiBus struct {
	fileDescriptor *os.File
}

// will be used if writes are done through ioctl instead of just using file descriptor
// Data that is passed to/from ioctl calls
// type spiIoctlData struct {

// }

// Constants used by ioctl
const (
	SPI_READ  = 1
	SPI_WRITE = 0

	SPI_IOC_WR_MODE = iota
	SPI_IOC_RD_MODE
	SPI_IOC_WR_BITS_PER_WORD
	SPI_IOC_RD_BITS_PER_WORD
	SPI_IOC_WR_MAX_SPEED_HZ
	SPI_IOC_RD_MAX_SPEED_HZ
)

func NewSpiBus(busFilePath string, busSpeedHz uint) *SpiBus {

	if 100 < busSpeedHz || busSpeedHz > 10000000 {
		log.Fatal("Bus speed is out of range", busSpeedHz)
	}

	file, err := os.OpenFile(busFilePath, os.O_RDWR, os.ModeExclusive)
	if err != nil {
		log.Fatal(err)
	}

	configBus(file, SPI_IOC_WR_MODE, 0)
	configBus(file, SPI_IOC_RD_MODE, 0)
	configBus(file, SPI_IOC_WR_BITS_PER_WORD, 8)
	configBus(file, SPI_IOC_RD_BITS_PER_WORD, 8)
	configBus(file, SPI_IOC_WR_MAX_SPEED_HZ, busSpeedHz)
	configBus(file, SPI_IOC_RD_MAX_SPEED_HZ, busSpeedHz)

	return &SpiBus{
		fileDescriptor: file,
	}
}

// Set the value
func configBus(file *os.File, command, value uint) {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, 3, uintptr(file.Fd()), uintptr(unsafe.Pointer(&command)), uintptr(unsafe.Pointer(&value)))
	if err != 0 {
		log.Fatal("Error attempting to configure SPI bus", err)
	}
}

// Write data to the bus
func (bus *SpiBus) Write(data []byte) (n int, err error) {

	n, err = bus.fileDescriptor.Write(data)
	if n != len(data) {
		log.Fatal("Failed to write all of the bytes, ", n, " instead of ", len(data))
	}
	bus.fileDescriptor.Sync() // flush data to be sure it's been written

	return
}
