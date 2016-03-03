package rpigpio

import "sync"

// PinFunction represents the BCM2835 function of a channel
type PinFunction uint8

// PinDirection represents the direction (IN/OUT) of a channel
type PinDirection PinFunction

// PinState represents the state of an output channel: high/low
type PinState uint8

// Pull up/down/off
type Pull uint8

// Status represents the status of the RpiGpio package
type Status uint8

// RpiGpio holds all data for a RPi GPIO implementation
type RpiGpio struct {
	lock        sync.Mutex
	mem         []uint32
	mem8        []uint8
	pinToBCMPin [27]int
	rpi         *RpiInfo
	status      Status
}
