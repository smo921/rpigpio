package rpigpio

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func rpigpioTestInit() *RpiGpio {
	gpio := new(RpiGpio)
	gpio.mem = make([]uint32, 54)
	return gpio
}

func TestDirection(t *testing.T) {
	gpio := rpigpioTestInit()
	var c uint8
	initRegisterValues := [5]uint32{
		0x00000000,
		0xFFFFFFFF,
		0xF0F0F0F0,
		0xFFFF0000,
		0x0000FFFF,
	}
	for c = 0; c < 54; c++ {
		fsel := c / 10
		for x := range initRegisterValues {
			gpio.mem[fsel] = initRegisterValues[x]
			err := testPinDirection(gpio, c)
			if err != nil {
				t.Error(fmt.Errorf("Init Register Value: %08X ; %s", initRegisterValues[x], err))
			}
		}
	}
}

func TestMmapFile(t *testing.T) {
	file, err := os.OpenFile("./test/gpiomem",
		os.O_CREATE|os.O_TRUNC|os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		t.Error(fmt.Errorf("Error opening ./test/gpiomem: %s", err))
	}
	defer os.Remove("./test/gpiomem")

	// Initialize file to memLength bytes
	var buf []byte
	buf = make([]byte, memLength)
	file.Write(buf)

	gpio := rpigpioTestInit()
	err = gpio.mmapFile(file)
	if err != nil {
		t.Error(fmt.Errorf("Error in mmapFile: %s", err))
	}
	gpio.mem[5] = 0x12345678
	gpio.Close()
	file.Close()

	file, err = os.OpenFile("./test/gpiomem", os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		t.Error(fmt.Errorf("Error opening ./test/gpiomem: %s", err))
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

func TestPull(t *testing.T) {
	gpio := rpigpioTestInit()
	err := gpio.setPull(PULLDOWN)
	if gpio.mem[pullUpDownOffset] != uint32(PULLDOWN) || err != nil {
		t.Error("setPull(DOWN) failed:", err)
	}
	err = gpio.setPull(PULLUP)
	if gpio.mem[pullUpDownOffset] != uint32(PULLUP) || err != nil {
		t.Error("setPull(DOWN) failed:", err)
	}
	err = gpio.setPull(PULLOFF)
	if gpio.mem[pullUpDownOffset] != uint32(PULLOFF) || err != nil {
		t.Error("setPull(DOWN) failed:", err)
	}
	err = gpio.setPull(5)
	if err == nil {
		t.Error("setPull(5) should return error:", err)
	}

}

func TestRead(t *testing.T) {
	gpio := rpigpioTestInit()
	var c uint8
	for c = 0; c < 54; c++ {
		register := (c / 32) + pinLevelOffset
		shift := c % 32
		// Set pin to HIGH
		gpio.mem[register] |= (1 << shift)
		val := gpio.Read(c)
		if val != HIGH {
			t.Error("Expected pin to be HIGH")
		}
		// clear all bits for pin ; ie set pin to LOW
		gpio.mem[register] &^= (gpioPinMask << shift)
		val = gpio.Read(c)
		if val != LOW {
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
	gpio := rpigpioTestInit()
	var c uint8
	for c = 0; c < 54; c++ {
		setRegister := (c / 32) + setOffset
		clearRegister := (c / 32) + clearOffset
		shift := c % 32

		gpio.mem[clearRegister] = 0
		gpio.Write(c, LOW)
		val := gpio.mem[clearRegister]
		if val != (1 << shift) {
			t.Error("Expected pin to be LOW")
		}

		gpio.mem[setRegister] = 0
		gpio.Write(c, HIGH)
		val = gpio.mem[setRegister]
		if val != (1 << shift) {
			t.Error("Expected pin to be HIGH")
		}
	}
}

func testPinDirection(gpio *RpiGpio, c uint8) (err error) {
	var mask uint32
	var val uint32

	fsel := c / 10
	shift := (c % 10) * 3
	mask = (gpioPinMask << shift)

	gpio.Direction(c, IN)
	val = gpio.mem[fsel] & mask
	if val != uint32(IN) {
		err = fmt.Errorf("Failed to set pin %d to input", c)
	}

	gpio.Direction(c, OUT)
	val = (gpio.mem[fsel] & mask) >> shift
	if val != uint32(OUT) {
		err = fmt.Errorf("Failed to set pin %d to output", c)
	}
	return
}
