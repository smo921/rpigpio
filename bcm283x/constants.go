package bcm283x

import "github.com/smo921/rpigpio/gpio"

const (
	// Only needed if we want to mmap /dev/mem
	bcm2835Base    = 0x20000000
	gpioBaseOffset = 0x200000
	piGPIOBase     = bcm2835Base + gpioBaseOffset
	bcm2835Max     = 0x20FFFFFF
	// ^^ can be removed if we do not want to support /dev/mem
	fselOffset          = 0
	setOffset           = 7
	clearOffset         = 10
	pinLevelOffset      = 13
	pullUpDownOffset    = 37
	pullUpDownClkOffset = 38

	memLength   = 4096
	gpioPinMask = 7
)

// Enumerate avaialable channel directions
const (
	IN  gpio.PinDirection = gpio.PinDirection(INPUT)
	OUT gpio.PinDirection = gpio.PinDirection(OUTPUT)
)

// Enumerate avaialable channel functions
const (
	INPUT  gpio.PinFunction = 0
	OUTPUT                  = 1
	ALT0                    = 4
	ALT1                    = 5
	ALT2                    = 6
	ALT3                    = 7
	ALT4                    = 3
	ALT5                    = 2
)

// Enumerate possible pin states
const (
	LOW gpio.PinState = iota
	HIGH
)

// Pull up/down/off
const (
	PULLOFF gpio.Pull = iota
	PULLDOWN
	PULLUP
)
