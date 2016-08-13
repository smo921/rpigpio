package bcm283x

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/smo921/rpigpio/gpio"
)

func bcmGpioTestInit() *BcmGpio {
	bcm := new(BcmGpio)
	bcm.mem = make([]uint32, 54)
	return bcm
}

func TestMmapFile(t *testing.T) {
	file, err := os.OpenFile("../test/gpiomem",
		os.O_CREATE|os.O_TRUNC|os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		t.Error(fmt.Errorf("Error opening ./test/gpiomem: %s", err))
	}
	defer os.Remove("../test/gpiomem")

	// Initialize file to memLength bytes
	var buf []byte
	buf = make([]byte, memLength)
	file.Write(buf)

	gpio := bcmGpioTestInit()
	err = gpio.mmapFile(file)
	if err != nil {
		t.Error(fmt.Errorf("Error in mmapFile: %s", err))
	}
	gpio.mem[5] = 0x12345678
	gpio.Close()
	file.Close()

	file, err = os.OpenFile("../test/gpiomem", os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		t.Error(fmt.Errorf("Error opening ../test/gpiomem: %s", err))
	}
	defer file.Close()
	defer gpio.Close() // close mmap before closing file

	err = gpio.mmapFile(file)
	if err != nil {
		t.Error(fmt.Errorf("Error in mmapFile: %s", err))
	}
	if gpio.mem[5] != 0x12345678 {
		t.Error("Data read from mmap does not match what was written")
	}
}

func TestDirection(t *testing.T) {
	bcm := bcmGpioTestInit()
	var p gpio.Pin
	initRegisterValues := [5]uint32{
		0x00000000,
		0xFFFFFFFF,
		0xF0F0F0F0,
		0xFFFF0000,
		0x0000FFFF,
	}
	for p = 0; p < 54; p++ {
		fsel := p / 10
		for x := range initRegisterValues {
			bcm.mem[fsel] = initRegisterValues[x]
			err := testPinDirection(bcm, p)
			if err != nil {
				t.Error(fmt.Errorf("Init Register Value: %08X ; %s", initRegisterValues[x], err))
			}
		}
	}
}

func TestPull(t *testing.T) {
	gpio := bcmGpioTestInit()
	err := gpio.setPull(PULLDOWN)
	if gpio.mem[pullUpDownOffset] != uint32(PULLDOWN) || err != nil {
		t.Error("setPull(DOWN) failed:", err)
	}
	err = gpio.setPull(PULLUP)
	if gpio.mem[pullUpDownOffset] != uint32(PULLUP) || err != nil {
		t.Error("setPull(UP) failed:", err)
	}
	/*
		err = gpio.setPull(PULLOFF)
		if gpio.mem[pullUpDownOffset] != uint32(PULLOFF) || err != nil {
			t.Error("setPull(OFF) failed:", err)
		}
	*/
	err = gpio.setPull(5)
	if err == nil {
		t.Error("setPull(5) should return error:", err)
	}
}

func TestRead(t *testing.T) {
	bcm := bcmGpioTestInit()
	var p gpio.Pin
	for p = 0; p < 54; p++ {
		register := (p / 32) + pinLevelOffset
		shift := p % 32
		// Set pin to HIGH
		bcm.mem[register] |= (1 << shift)
		val := bcm.Read(p)
		if val != gpio.HIGH {
			t.Error("Expected pin to be HIGH")
		}
		// clear all bits for pin ; ie set pin to LOW
		bcm.mem[register] &^= (gpioPinMask << shift)
		val = bcm.Read(p)
		if val != gpio.LOW {
			t.Error("Expected pin to be LOW")
		}
	}
}

func TestShortWait(t *testing.T) {
	var x uint32
	for x = 150; x < 10000; x *= 2 {
		start := time.Now()
		shortWait(x)
		dur := time.Since(start).Seconds()
		fmt.Printf("shortWait(%d) took: %f sec (%0.2f hz)\n", x, dur, float64(x)/dur)
	}
}

func TestWrite(t *testing.T) {
	bcm := bcmGpioTestInit()
	var p gpio.Pin
	for p = 0; p < 54; p++ {
		setRegister := (p / 32) + setOffset
		clearRegister := (p / 32) + clearOffset
		shift := p % 32

		bcm.mem[clearRegister] = 0
		bcm.Write(p, gpio.LOW)
		val := bcm.mem[clearRegister]
		if val != (1 << shift) {
			t.Error("Expected pin to be LOW")
		}

		bcm.mem[setRegister] = 0
		bcm.Write(p, gpio.HIGH)
		val = bcm.mem[setRegister]
		if val != (1 << shift) {
			t.Error("Expected pin to be HIGH")
		}
	}
}

func testPinDirection(bcm *BcmGpio, p gpio.Pin) (err error) {
	var mask uint32
	var val uint32

	fsel := p / 10
	shift := (p % 10) * 3
	mask = (gpioPinMask << shift)

	bcm.Direction(p, IN)
	val = bcm.mem[fsel] & mask
	if val != uint32(IN) {
		err = fmt.Errorf("Failed to set pin %d to input", p)
	}

	bcm.Direction(p, OUT)
	val = (bcm.mem[fsel] & mask) >> shift
	if val != uint32(OUT) {
		err = fmt.Errorf("Failed to set pin %d to output", p)
	}
	return
}
