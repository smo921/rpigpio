package rpigpio

import (
	"errors"
	"fmt"

	"github.com/smo921/rpigpio/bcm283x"
	"github.com/smo921/rpigpio/gpio"
)

var piPinToBCMPinRev2 = [27]int8{
	-1, -1, -1, 2, -1, 3, -1, 4, 14, -1, 15, 17, 18, 27, -1,
	22, 23, -1, 24, 10, -1, 9, 25, 11, 8, -1, 7,
}

// Deteriming GPIO number
func (rpi *RpiGpio) getBCMGpio(pin gpio.Pin) (gpio.Pin, error) {
	if rpi.mode == GPIO {
		return pin, nil
	} else if rpi.mode == PI && int(pin) > len(rpi.pinToBCMPin) {
		return 255, fmt.Errorf("Pin %d must be < %d", pin, len(rpi.pinToBCMPin))
	} else if rpi.pinToBCMPin[pin] == -1 {
		return 255, fmt.Errorf("Pin %d not valid for GPIO", pin)
	}
	return gpio.Pin(rpi.pinToBCMPin[pin]), nil
}

// NewGPIO sets up a new GPIO object
func NewGPIO() (*RpiGpio, error) {
	var err error
	rpi := new(RpiGpio)
	rpi.bcm, err = bcm283x.New()
	if err != nil {
		return nil, err
	}
	rpi.mode = GPIO
	rpi.status = NEW
	rpi.info = new(RpiInfo)
	err = rpi.info.GetCPUInfo()
	if err != nil {
		return nil, err
	}
	switch int(rpi.info.piRevision) {
	case 2:
		rpi.pinToBCMPin = piPinToBCMPinRev2
	default:
		return nil, errors.New("Unknown Raspberry Pi hardware")
	}
	// set gpio "direction"  (in/out??)
	// pinTopin = pinToGpiopinRev??
	rpi.status = OK
	return rpi, nil
}

// Direction sets the pin as either input or output
func (rpi *RpiGpio) Direction(p gpio.Pin, d gpio.PinDirection) error {
	pin, err := rpi.getBCMGpio(p)
	if err != nil {
		return err
	}
	return rpi.bcm.Direction(pin, d)
}

// Mode sets the pin interpretation for the rpigpio functions
func (rpi *RpiGpio) Mode(m Mode) error {
	if m != GPIO && m != PI {
		return fmt.Errorf("Mode must be GPIO or PI")
	}
	rpi.mode = m
	return nil
}

// Pull sets the direction of the built-in pull-up/pull-down resistor
func (rpi *RpiGpio) Pull(p gpio.Pin, d gpio.Pull) error {
	pin, err := rpi.getBCMGpio(p)
	if err != nil {
		return err
	}
	return rpi.bcm.Pull(pin, d)
}

func (rpi *RpiGpio) Read(p gpio.Pin) (gpio.PinState, error) {
	pin, err := rpi.getBCMGpio(p)
	if err != nil {
		return 255, err
	}
	return rpi.bcm.Read(pin), nil
}

func (rpi *RpiGpio) Write(p gpio.Pin, s gpio.PinState) error {
	pin, err := rpi.getBCMGpio(p)
	if err != nil {
		return err
	}
	rpi.bcm.Write(pin, s)
	return nil
}

// Cleanup the pin ; reset to INPUT and pull up/down to off
func (rpi *RpiGpio) Cleanup(pin uint8) {
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
