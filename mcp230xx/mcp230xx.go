package mcp230xx

import "github.com/davecheney/i2c"

const (
	iodir = 0x00
	gpio  = 0x12
	gppu  = 0x0C
)

// MCP230xx represents a MCP230xx series GPIO extender
type MCP230xx struct {
	i2c               *i2c.I2C
	iodir, gppu, gpio []byte
}

// New initializes a MCP230xx
func New(address uint8, bus, gpioPins int) (*MCP230xx, error) {
	i2c, err := i2c.New(address, bus)
	if err != nil {
		return nil, err
	}
	mcp := new(MCP230xx)
	mcp.i2c = i2c
	p := gpioPins / 8
	mcp.iodir = make([]byte, p)
	mcp.gppu = make([]byte, p)
	mcp.gpio = make([]byte, p)
	return mcp, nil
}

// Direction sets Pin to PinDirection (IN/OUT)
func (mcp *MCP230xx) Direction(p Pin, d PinDirection) error {

}

// Pull sets the pull direction on Pin p (Up/Down)
func (mcp *MCP230xx) Pull(p Pin, d Pull) error {

}

// Read gets the current value (0/1) from Pin p
func (mcp *MCP230xx) Read(p Pin) error {

}

// Write sets the state of Pin p (0/1)
func (mcp *MCP230xx) Write(p Pin, s PinState) error {

}

// WriteGPIO writes the byte value to the GPIO register
func (mcp *MCP230xx) writeGPIO(b []byte) error {
	_, e := mcp.write(gpio, b)
	return e
}

// WriteIODIR writes the byte value to the IODIR register
func (mcp *MCP230xx) writeIODIR(b []byte) error {
	_, e := mcp.write(iodir, b)
	return e
}

// WriteGPPU writes the byte value to the GPIO Pull Up/Down register
func (mcp *MCP230xx) writeGPPU(b []byte) error {
	_, e := mcp.write(gppu, b)
	return e
}

func (mcp *MCP230xx) write(register int, val []byte) (int, error) {
	data := []byte{byte(len(val) & 0xFF)}
	data = append(data, val...)
	return mcp.i2c.Write(data)
}
