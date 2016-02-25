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
