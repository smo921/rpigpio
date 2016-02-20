package rpigpio

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"syscall"
	"unsafe"
)

// ChannelDirection represents the direction (IN/OUT) of a channel
type ChannelDirection uint8

// ChannelState represents the state of an output channel: high/low
type ChannelState uint8

// Pull up/down/off
type Pull uint8

// Status represents the status of the RpiGpio package
type Status uint8

const (
	// Only needed if we want to mmap /dev/mem
	bcm2835Base    = 0x20000000
	gpioBaseOffset = 0x200000
	piGPIOBase     = bcm2835Base + gpioBaseOffset
	bcm2835Max     = 0x20FFFFFF
	// ^^ can be removed if we do not want to support /dev/mem
	fselOffset          = 0
	setOffset           = 7
	clearOffset         = 10
	pinLevelOffset      = 13
	pullUpDownOffset    = 37
	pullUpDownClkOffset = 38

	memLength   = 4096
	gpioPinMask = 7
)

// Enumerate avaialable channel directions
const (
	IN ChannelDirection = iota
	OUT
)

// Enumerate possible channel states
const (
	LOW ChannelState = iota
	HIGH
)

// Pull up/down/off
const (
	PULLOFF Pull = iota
	PULLDOWN
	PULLUP
)

// Status of the RpiGpio package
const (
	NEW Status = iota
	OK
	ERR
)

// RpiGpio holds all data for a RPi GPIO implementation
type RpiGpio struct {
	lock         sync.Mutex
	mem          []uint32
	mem8         []uint8
	pinToChannel [27]int
	rpi          *RpiInfo
	status       Status
}

var pinToGpioChannelRev2 = [27]int{
	-1, -1, -1, 2, -1, 3, -1, 4, 14, -1, 15, 17, 18, 27, -1,
	22, 23, -1, 24, 10, -1, 9, 25, 11, 8, -1, 7,
}

// Deteriming GPIO number
func getGpioNumber(channel int) (uint, error) {
	return 0, nil
}

// Close cleans up the RpiGpio resources
func (gpio *RpiGpio) Close() error {
	//event_cleanup_all()???  When we implement an event handler??
	gpio.lock.Lock()
	defer gpio.lock.Unlock()
	return syscall.Munmap(gpio.mem8)
}

// NewGPIO sets up a new GPIO object
func NewGPIO() (*RpiGpio, error) {
	var err error
	gpio := new(RpiGpio)
	gpio.status = NEW
	gpio.rpi = new(RpiInfo)
	gpio.rpi.GetCPUInfo()
	switch int(gpio.rpi.piRevision) {
	case 2:
		gpio.pinToChannel = pinToGpioChannelRev2
	default:
		return nil, errors.New("Unknown Raspberry Pi hardware")
	}
	// set gpio "direction"  (in/out??)
	// pinToChannel = pinToGpioChannelRev??
	err = gpio.openGPIO()
	if err != nil {
		return nil, err
	}
	gpio.status = OK
	return gpio, nil
}

// Direction configures the direction (IN/OUT) of the channel
func (gpio *RpiGpio) Direction(channel uint8, direction ChannelDirection) (err error) {
	// Check package status is OK
	// do some error checking ; verify channel and direction are valid, etc
	// call c_gpio::setup_one()
	fsel := channel / 10
	shift := (channel % 10) * 3
	switch direction {
	case IN:
		gpio.mem[fsel] = gpio.mem[fsel] &^ (gpioPinMask << shift)
	case OUT:
		gpio.mem[fsel] = (gpio.mem[fsel] &^ (gpioPinMask << shift)) | (1 << shift)
	default:
		errString := fmt.Sprintf("Unknown channel direction: %d", direction)
		err = errors.New(errString)
	}
	return
}

// Pull sets or clears the internal pull up/down resistor for a GPIO channel
func (gpio *RpiGpio) Pull(channel uint8, direction Pull) error {
	clkRegister := (channel / 32) + pullUpDownClkOffset
	shift := channel % 32

	switch direction {
	case PULLOFF:
		gpio.mem[pullUpDownOffset] &^= 3
	case PULLDOWN, PULLUP:
		gpio.mem[pullUpDownOffset] = (gpio.mem[pullUpDownOffset] &^ 3) | uint32(direction)
	default:
		errString := fmt.Sprintf("Unknown pull direction: %d", direction)
		return errors.New(errString)
	}

	shortWait(150)
	gpio.mem[clkRegister] = 1 << shift
	shortWait(150)
	gpio.mem[pullUpDownOffset] &^= 3
	gpio.mem[clkRegister] = 0
	return nil
}

// Read value from channel
func (gpio *RpiGpio) Read(channel uint8) ChannelState {
	pinLevelRegister := (channel / 32) + pinLevelOffset
	shift := channel % 32
	if gpio.mem[pinLevelRegister]&(1<<shift) != 0 {
		return HIGH
	}
	return LOW
}

// Write value (0/1) to channel
func (gpio *RpiGpio) Write(channel uint8, state ChannelState) error {
	reg := channel / 32
	shift := channel % 32
	gpio.lock.Lock()

	if state == HIGH {
		reg += setOffset
	} else if state == LOW {
		reg += clearOffset
	} else {
		err := fmt.Sprintf("Unknown channel state: %d", state)
		return errors.New(err)
	}
	gpio.mem[reg] = 1 << shift
	gpio.lock.Unlock()
	return nil
}

// Cleanup the channel ; reset to INPUT and pull up/down to off
func (gpio *RpiGpio) Cleanup(channel uint8) {
	// Verify channel is valid, package status is OK, etc
	// get gpio number from channel
	// call c_gpio::cleanup_one()
	//    * call event_cleanup()
	//    * set gpio_direction = -1
	//    * set gpio to INPUT and pull up/down to off
	//    * set found for error checking later on (if working on > 1 channel at a time)

}

// eventCleanup
func eventCleanup(gpio uint) {
	// event_gpio.c:403
}

func (gpio *RpiGpio) openGPIO() (err error) {
	file, err := os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		fmt.Println("Error opening /dev/gpiomem: ", err)
		return
	}
	defer file.Close()
	return gpio.mmapFile(file)
}

func (gpio *RpiGpio) mmapFile(f *os.File) (err error) {
	gpio.lock.Lock()
	defer gpio.lock.Unlock()
	// Memory map GPIO registers to byte array
	gpio.mem8, err = syscall.Mmap(
		int(f.Fd()),
		0,
		memLength,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED)
	if err != nil {
		return
	}

	// Convert 8-bit slice to 32-bit slice
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&gpio.mem8))
	header.Len /= 4
	header.Cap /= 4
	gpio.mem = *(*[]uint32)(unsafe.Pointer(&header))
	return nil
}

// run cnt nop operations
func shortWait(cnt uint32)
