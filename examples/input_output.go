package main

import (
	"fmt"
	"time"

	"github.com/smo921/rpigpio"
	"github.com/smo921/rpigpio/gpio"
)

func main() {
	pi, _ := rpigpio.NewGPIO()
	pi.Direction(7, rpigpio.IN)
	pi.Direction(8, rpigpio.OUT)  // green
	pi.Direction(25, rpigpio.OUT) // red
	pi.Pull(7, rpigpio.PULLOFF)

	pi.Write(8, gpio.LOW)
	pi.Write(25, gpio.HIGH)
	var val gpio.PinState
	for n := 0; n < 10; n++ {
		val, _ = pi.Read(7)
		fmt.Println("Read:", val)
		if val == 1 {
			pi.Write(25, gpio.LOW)
			pi.Write(8, gpio.HIGH)
		} else {
			pi.Write(25, gpio.HIGH)
			pi.Write(8, gpio.LOW)
		}
		time.Sleep(time.Second)
	}
	pi.Write(25, gpio.LOW)
	pi.Write(8, gpio.LOW)
}
