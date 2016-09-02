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
	pi.Pull(7, rpigpio.PULLUP)
	var val gpio.PinState
	for n := 0; n < 100; n++ {
		val, _ = pi.Read(7)
		fmt.Println("Read:", val)
		time.Sleep(time.Second)
	}
}
