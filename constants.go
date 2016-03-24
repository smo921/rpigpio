package rpigpio

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

// How pin numbers are interpreted (pi header vs bcm gpio)
const (
	GPIO Mode = iota
	PI
)

// Enumerate avaialable channel directions
const (
	IN  PinDirection = PinDirection(INPUT)
	OUT PinDirection = PinDirection(OUTPUT)
)

// Enumerate avaialable channel functions
const (
	INPUT  PinFunction = 0
	OUTPUT             = 1
	ALT0               = 4
	ALT1               = 5
	ALT2               = 6
	ALT3               = 7
	ALT4               = 3
	ALT5               = 2
)

// Enumerate possible pin states
const (
	LOW PinState = iota
	HIGH
)

// Pull up/down/off
const (
	PULLOFF Pull = iota
	PULLDOWN
	PULLUP
)

// Status of the RpiGpio package
const (
	NEW Status = iota
	OK
	ERR
)
