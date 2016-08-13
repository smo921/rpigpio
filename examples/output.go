package main

import (
	"time"

	"github.com/smo921/rpigpio"
	"github.com/smo921/rpigpio/gpio"
)

func main() {
	pi, _ := rpigpio.NewGPIO()
	pi.Direction(7, rpigpio.OUT)

	for x := 0; x < 10; x++ {
		pi.Write(7, gpio.HIGH)
		time.Sleep(time.Second)
		pi.Write(7, gpio.LOW)
		time.Sleep(time.Second)
	}

}
