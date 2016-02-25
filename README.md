# RpiGpio
RpiGpio is a Go package for interacting with the GPIO pins.  This package is heavily
influenced by the Python RPi.GPIO module and Stian Eilkland's go-rpio package.

# Example Code
This example demonstrates how to toggle an output pin from high to low

```
package main

import (
	"time"

	"github.com/smo921/rpigpio"
)

func main() {
	gpio, _ := rpigpio.NewGPIO()
	gpio.Direction(7, rpigpio.OUT)

	for x := 0; x < 10; x++ {
		gpio.Write(7, rpigpio.HIGH)
		time.Sleep(time.Second)
		gpio.Write(7, rpigpio.LOW)
		time.Sleep(time.Second)
	}

}
```

# References
## Raspberry Pi hardware
* https://www.raspberrypi.org/wp-content/uploads/2012/02/BCM2835-ARM-Peripherals.pdf
* http://elinux.org/RPi_HardwareHistory

## GPIO
* http://sourceforge.net/p/raspberry-gpio-python/wiki/BasicUsage/
* https://blog.eikeland.se/2013/07/30/go-gpio-library-for-raspberry-pi/
* https://github.com/stianeikeland/go-rpio
* https://pypi.python.org/pypi/RPi.GPIO
* http://abyz.co.uk/rpi/pigpio/examples.html#Misc_code
* https://learn.sparkfun.com/tutorials/raspberry-gpio/python-rpigpio-api
* http://raspi.tv/2013/rpi-gpio-basics-6-using-inputs-and-outputs-together-with-rpi-gpio-pull-ups-and-pull-downs
