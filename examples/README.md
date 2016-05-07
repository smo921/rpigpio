# RpiGpio Examples

## Output
Very basic example showing how to toggle an LED on/off

[Example Code](./output.go)
![Output Breadboard](./output_bb.png)

## Input
Demonstration of reading input from a push button.  This does not use the pull up/down feature of the Raspberry Pi GPIO pins.  Instead the pull down is built into the circuit itself.

[Example Code](./input.go)
![Input Breadboard](./input_bb.png)

A second [example](./input2.go) sets the internal pull up/down resistor state instead of adding one
to the circuit on the breadboard.

## Input/Output
Building on the Input and Output examples above we toggle between a red and green led based on the switch position.

[Example Code](./input_output.go)
![I/O Breadboard](./input_output_bb.png)
