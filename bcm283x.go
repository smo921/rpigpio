package rpigpio

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

// Close cleans up the RpiGpio resources
func (gpio *RpiGpio) Close() error {
	//event_cleanup_all()???  When we implement an event handler??
	gpio.lock.Lock()
	defer gpio.lock.Unlock()
	return syscall.Munmap(gpio.mem8)
}

// Direction configures the direction (IN/OUT) of the pin
func (gpio *RpiGpio) Direction(pin uint8, direction PinDirection) (err error) {
	// Check package status is OK
	// do some error checking ; verify pin and direction are valid, etc
	// call c_gpio::setup_one()
	fsel := pin / 10
	shift := (pin % 10) * 3
	switch direction {
	case IN:
		gpio.mem[fsel] = gpio.mem[fsel] &^ (gpioPinMask << shift)
	case OUT:
		gpio.mem[fsel] = (gpio.mem[fsel] &^ (gpioPinMask << shift)) | (1 << shift)
	default:
		errString := fmt.Sprintf("Unknown pin direction: %d", direction)
		err = errors.New(errString)
	}
	return
}

// Pull sets or clears the internal pull up/down resistor for a GPIO pin
func (gpio *RpiGpio) Pull(pin uint8, direction Pull) error {
	clkRegister := (pin / 32) + pullUpDownClkOffset
	shift := pin % 32

	if err := gpio.setPull(direction); err != nil {
		return err
	}

	shortWait(150)
	gpio.mem[clkRegister] = 1 << shift
	shortWait(150)
	gpio.mem[pullUpDownOffset] &^= 3
	gpio.mem[clkRegister] = 0
	return nil
}

// Read value from pin
func (gpio *RpiGpio) Read(pin uint8) PinState {
	pinLevelRegister := (pin / 32) + pinLevelOffset
	shift := pin % 32
	if gpio.mem[pinLevelRegister]&(1<<shift) != 0 {
		return HIGH
	}
	return LOW
}

// Write value (0/1) to pin
func (gpio *RpiGpio) Write(pin uint8, state PinState) error {
	reg := pin / 32
	shift := pin % 32
	gpio.lock.Lock()

	if state == HIGH {
		reg += setOffset
	} else if state == LOW {
		reg += clearOffset
	} else {
		err := fmt.Sprintf("Unknown pin state: %d", state)
		return errors.New(err)
	}
	gpio.mem[reg] = 1 << shift
	gpio.lock.Unlock()
	return nil
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

func (gpio *RpiGpio) setPull(d Pull) error {
	switch d {
	case PULLOFF:
		gpio.mem[pullUpDownOffset] &^= 3
	case PULLDOWN, PULLUP:
		gpio.mem[pullUpDownOffset] = (gpio.mem[pullUpDownOffset] &^ 3) | uint32(d)
	default:
		errString := fmt.Sprintf("Unknown pull direction: %d", d)
		return errors.New(errString)
	}
	return nil
}

// run cnt nop operations
func shortWait(cnt uint32)
