package rpigpio

import "testing"

func rpigpioTestInit() *RpiGpio {
	gpio := new(RpiGpio)
	return gpio
}

func TestMode(t *testing.T) {
	gpio := rpigpioTestInit()
	if err := gpio.Mode(GPIO); err != nil {
		t.Error("Set rpigpio mode GPIO failed:", err)
	}
	if err := gpio.Mode(PI); err != nil {
		t.Error("Set rpigpio mode PI failed:", err)
	}
	if err := gpio.Mode(255); err == nil {
		t.Error("Set rpigpio mode 255 passed")
	}

}

func TestGetBCMGpio(t *testing.T) {
	gpio := rpigpioTestInit()
	gpio.Mode(PI) // In PI mode we should convert pin values
	gpio.pinToBCMPin = piPinToBCMPinRev2
	_, err := gpio.getBCMGpio(0)
	if err == nil {
		t.Error("Pin 0 should return error")
	}
	pin, err := gpio.getBCMGpio(3)
	if pin != 2 {
		t.Error("Pin 3 should return BCM gpio 2")
	}
	pin, err = gpio.getBCMGpio(16)
	if pin != 23 {
		t.Error("Pin 16 should return BCM gpio 23")
	}
}
