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
	// these values were determined experimentally by running spidev_test.c on a raspberry pi
	// they may not work on other systems
	SPI_IOC_WR_MODE = 0x40016B01
	SPI_IOC_RD_MODE = 0x80016B01
	SPI_IOC_WR_BITS_PER_WORD = 0x40016B03
	SPI_IOC_RD_BITS_PER_WORD = 0x80016B03
	SPI_IOC_WR_MAX_SPEED_HZ = 0x40046B04
	SPI_IOC_RD_MAX_SPEED_HZ = 0x80046B04
)

func NewSpiBus(busFilePath string, busSpeedHz uint) *SpiBus {

	if busSpeedHz < 100 || 10000000 < busSpeedHz {
		log.Fatal("Bus speed is out of range", busSpeedHz)
	}

	file, err := os.OpenFile(busFilePath, os.O_RDWR, os.ModeExclusive)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("WR_MODE")
	configBus(file, SPI_IOC_WR_MODE, 0)
	log.Print("RD_MODE")
	configBus(file, SPI_IOC_RD_MODE, 0)
	log.Print("WR_BITS")
	configBus(file, SPI_IOC_WR_BITS_PER_WORD, 8)
	log.Print("RD_BITS")
	configBus(file, SPI_IOC_RD_BITS_PER_WORD, 8)
	log.Print("WR_SPEED")
	configBus(file, SPI_IOC_WR_MAX_SPEED_HZ, busSpeedHz)
	log.Print("RD_SPEED")
	configBus(file, SPI_IOC_RD_MAX_SPEED_HZ, busSpeedHz)

	return &SpiBus{
		fileDescriptor: file,
	}
}

// Set the value
func configBus(file *os.File, command, value uint) {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(file.Fd()), uintptr(command), uintptr(unsafe.Pointer(&value)))
	if err != 0 {
		log.Fatal("Error attempting to configure SPI bus ", err)
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
