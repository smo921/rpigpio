package rpigpio

func rpigpioTestInit() *RpiGpio {
	gpio := new(RpiGpio)
	gpio.mem = make([]uint32, 54)
	return gpio
}
