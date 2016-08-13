package bcm283x

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"syscall"
	"unsafe"

	"github.com/smo921/rpigpio/gpio"
)

// BcmGpio holds data for interacting with the BCM283x SoC
type BcmGpio struct {
	lock sync.Mutex
	max  gpio.Pin
	mem  []uint32
	mem8 []uint8
}

// New creates a bcmGpio data structure
func New() (*BcmGpio, error) {
	bcm := new(BcmGpio)
	bcm.max = 53 // never have pin value > max
	err := bcm.open()
	if err != nil {
		return nil, err
	}
	return bcm, nil
}

// Close cleans up the bcmGpio resources
func (bcm *BcmGpio) Close() error {
	//event_cleanup_all()???  When we implement an event handler??
	bcm.lock.Lock()
	defer bcm.lock.Unlock()
	return syscall.Munmap(bcm.mem8)
}

// Direction configures the direction (IN/OUT) of the pin
func (bcm *BcmGpio) Direction(pin gpio.Pin, direction gpio.PinDirection) (err error) {
	// Check package status is OK
	// do some error checking ; verify pin and direction are valid, etc
	// call c_gpio::setup_one()
	fsel := pin / 10
	shift := (pin % 10) * 3
	switch direction {
	case IN:
		bcm.mem[fsel] = bcm.mem[fsel] &^ (gpioPinMask << shift)
	case OUT:
		bcm.mem[fsel] = (bcm.mem[fsel] &^ (gpioPinMask << shift)) | (1 << shift)
	default:
		errString := fmt.Sprintf("Unknown pin direction: %d", direction)
		err = errors.New(errString)
	}
	return
}

// Pull sets or clears the internal pull up/down resistor for a GPIO pin
func (bcm *BcmGpio) Pull(pin gpio.Pin, direction gpio.Pull) error {
	clkRegister := (pin / 32) + pullUpDownClkOffset
	shift := pin % 32

	if err := bcm.setPull(direction); err != nil {
		return err
	}

	shortWait(150)
	bcm.mem[clkRegister] = 1 << shift
	shortWait(150)
	bcm.mem[pullUpDownOffset] &^= 3
	bcm.mem[clkRegister] = 0
	return nil
}

// Read value from pin
func (bcm *BcmGpio) Read(pin gpio.Pin) gpio.PinState {
	pinLevelRegister := (pin / 32) + pinLevelOffset
	shift := pin % 32
	if bcm.mem[pinLevelRegister]&(1<<shift) != 0 {
		return gpio.HIGH
	}
	return gpio.LOW
}

// Write value (0/1) to pin
func (bcm *BcmGpio) Write(pin gpio.Pin, state gpio.PinState) error {
	reg := pin / 32
	shift := pin % 32
	bcm.lock.Lock()

	if state == gpio.HIGH {
		reg += setOffset
	} else if state == gpio.LOW {
		reg += clearOffset
	} else {
		err := fmt.Sprintf("Unknown pin state: %d", state)
		return errors.New(err)
	}
	bcm.mem[reg] = 1 << shift
	bcm.lock.Unlock()
	return nil
}

func (bcm *BcmGpio) open() (err error) {
	file, err := os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		fmt.Println("Error opening /dev/gpiomem: ", err)
		return
	}
	defer file.Close()
	return bcm.mmapFile(file)
}

func (bcm *BcmGpio) mmapFile(f *os.File) (err error) {
	bcm.lock.Lock()
	defer bcm.lock.Unlock()
	// Memory map GPIO registers to byte array
	bcm.mem8, err = syscall.Mmap(
		int(f.Fd()),
		0,
		memLength,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED)
	if err != nil {
		return
	}

	// Convert 8-bit slice to 32-bit slice
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&bcm.mem8))
	header.Len /= 4
	header.Cap /= 4
	bcm.mem = *(*[]uint32)(unsafe.Pointer(&header))
	return nil
}

func (bcm *BcmGpio) setPull(d gpio.Pull) error {
	switch d {
	case PULLOFF:
		bcm.mem[pullUpDownOffset] &^= 3
	case PULLDOWN, PULLUP:
		bcm.mem[pullUpDownOffset] = (bcm.mem[pullUpDownOffset] &^ 3) | uint32(d)
	default:
		errString := fmt.Sprintf("Unknown pull direction: %d", d)
		return errors.New(errString)
	}
	return nil
}

// run cnt nop operations
func shortWait(cnt uint32)
