package main

import (
	"fmt"
	"time"

	"github.com/smo921/rpigpio"
)

func main() {
	gpio, _ := rpigpio.NewGPIO()
	gpio.Direction(7, rpigpio.IN)
	gpio.Direction(8, rpigpio.OUT)  // green
	gpio.Direction(25, rpigpio.OUT) // red
	gpio.Pull(7, rpigpio.PULLOFF)

	gpio.Write(8, rpigpio.LOW)
	gpio.Write(25, rpigpio.HIGH)
	var val rpigpio.PinState
	for n := 0; n < 10; n++ {
		val, _ = gpio.Read(7)
		fmt.Println("Read:", val)
		if val == 1 {
			gpio.Write(25, rpigpio.LOW)
			gpio.Write(8, rpigpio.HIGH)
		} else {
			gpio.Write(25, rpigpio.HIGH)
			gpio.Write(8, rpigpio.LOW)
		}
		time.Sleep(time.Second)
	}
	gpio.Write(25, rpigpio.LOW)
	gpio.Write(8, rpigpio.LOW)
}
