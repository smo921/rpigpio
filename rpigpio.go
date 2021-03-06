package rpigpio

import (
	"errors"
	"fmt"
)

var piPinToBCMPinRev2 = [27]int8{
	-1, -1, -1, 2, -1, 3, -1, 4, 14, -1, 15, 17, 18, 27, -1,
	22, 23, -1, 24, 10, -1, 9, 25, 11, 8, -1, 7,
}

// Deteriming GPIO number
func (gpio *RpiGpio) getBCMGpio(pin Pin) (Pin, error) {
	if gpio.mode == GPIO {
		return pin, nil
	} else if gpio.mode == PI && int(pin) > len(gpio.pinToBCMPin) {
		return 255, fmt.Errorf("Pin %d must be < %d", pin, len(gpio.pinToBCMPin))
	} else if gpio.pinToBCMPin[pin] == -1 {
		return 255, fmt.Errorf("Pin %d not valid for GPIO", pin)
	}
	return Pin(gpio.pinToBCMPin[pin]), nil
}

// NewGPIO sets up a new GPIO object
func NewGPIO() (*RpiGpio, error) {
	var err error
	gpio := new(RpiGpio)
	gpio.bcm = NewBCMGPIO()
	gpio.mode = GPIO
	gpio.status = NEW
	gpio.rpi = new(RpiInfo)
	gpio.rpi.GetCPUInfo()
	switch int(gpio.rpi.piRevision) {
	case 2:
		gpio.pinToBCMPin = piPinToBCMPinRev2
	default:
		return nil, errors.New("Unknown Raspberry Pi hardware")
	}
	// set gpio "direction"  (in/out??)
	// pinTopin = pinToGpiopinRev??
	err = gpio.bcm.open()
	if err != nil {
		return nil, err
	}
	gpio.status = OK
	return gpio, nil
}

// Direction sets the pin as either input or output
func (gpio *RpiGpio) Direction(p Pin, d PinDirection) error {
	pin, err := gpio.getBCMGpio(p)
	if err != nil {
		return err
	}
	return gpio.bcm.Direction(pin, d)
}

// Mode sets the pin interpretation for the rpigpio functions
func (gpio *RpiGpio) Mode(m Mode) error {
	if m != GPIO && m != PI {
		return fmt.Errorf("Mode must be GPIO or PI")
	}
	gpio.mode = m
	return nil
}

// Pull sets the direction of the built-in pull-up/pull-down resistor
func (gpio *RpiGpio) Pull(p Pin, d Pull) error {
	pin, err := gpio.getBCMGpio(p)
	if err != nil {
		return err
	}
	return gpio.bcm.Pull(pin, d)
}

func (gpio *RpiGpio) Read(p Pin) (PinState, error) {
	pin, err := gpio.getBCMGpio(p)
	if err != nil {
		return 255, err
	}
	return gpio.bcm.Read(pin), nil
}

func (gpio *RpiGpio) Write(p Pin, s PinState) error {
	pin, err := gpio.getBCMGpio(p)
	if err != nil {
		return err
	}
	gpio.bcm.Write(pin, s)
	return nil
}

// Cleanup the pin ; reset to INPUT and pull up/down to off
func (gpio *RpiGpio) Cleanup(pin uint8) {
	// Verify pin is valid, package status is OK, etc
	// get gpio number from pin
	// call c_gpio::cleanup_one()
	//    * call event_cleanup()
	//    * set gpio_direction = -1
	//    * set gpio to INPUT and pull up/down to off
	//    * set found for error checking later on (if working on > 1 pin at a time)
}

// eventCleanup
func eventCleanup(gpio uint) {
	// event_gpio.c:403
}
