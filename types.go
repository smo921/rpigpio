package rpigpio

import "github.com/smo921/rpigpio/bcm283x"

// Mode pi or bcm
type Mode uint8

// Status represents the status of the RpiGpio package
type Status uint8

// RpiGpio holds all data for a RPi GPIO implementation
type RpiGpio struct {
	bcm         *bcm283x.BcmGpio
	pinToBCMPin [27]int8
	mode        Mode
	info        *RpiInfo
	status      Status
}

// How pin numbers are interpreted (pi header vs bcm gpio)
const (
	GPIO Mode = iota
	PI
)

// Status of the RpiGpio package
const (
	NEW Status = iota
	OK
	ERR
)
