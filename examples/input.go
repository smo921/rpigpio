package main

import (
	"fmt"
	"time"
	"github.com/smo921/rpigpio"
)

func main() {
	gpio, _ := rpigpio.NewGPIO()
	gpio.Direction(7, rpigpio.IN)
	gpio.Pull(7, rpigpio.PULLOFF)
	var val rpigpio.PinState
	for n:=0; n< 10; n++ {
		val, _ = gpio.Read(7)
		fmt.Println("Read:", val)
		time.Sleep(time.Second)
	}
}
