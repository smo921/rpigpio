package gpio

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

// Enumerate possible pin states
const (
	LOW  PinState = 0
	HIGH PinState = 1
)
