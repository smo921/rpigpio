package rpigpio

import "sync"

// Mode pi or bcm
type Mode uint8

// Pin is a gpio or pi header pin number
type Pin uint8

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

// bcmGpio holds data for interacting with the BCM283x SoC
type bcmGpio struct {
	lock sync.Mutex
	mem  []uint32
	mem8 []uint8
}

// RpiGpio holds all data for a RPi GPIO implementation
type RpiGpio struct {
	bcm         *bcmGpio
	pinToBCMPin [27]int8
	mode        Mode
	rpi         *RpiInfo
	status      Status
}
