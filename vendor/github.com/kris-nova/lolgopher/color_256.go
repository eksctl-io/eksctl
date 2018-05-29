package lol

import (
	"fmt"
	"math"
)

type color256 struct {
	R, G, B uint8
}

func (c *color256) rainbow(freq float64, i float64) {
	c.R = uint8(math.Floor(math.Sin(freq*i+0)*127)) + 128
	c.G = uint8(math.Floor(math.Sin(freq*i+2.0*math.Pi/3.0)*127)) + 128
	c.B = uint8(math.Floor(math.Sin(freq*i+4.0*math.Pi/3.0)*127)) + 128
}

func (c *color256) format() []byte {

	// Support for grayscale colors
	if c.R == c.G && c.G == c.B && int(c.R) > 232 {
		return []byte(fmt.Sprintf("\x1b[38;5;%dm", int(c.R)))
	}

	// Math for 256: We use the closest value possible
	r6 := (uint16(c.R) * 3) / 255
	g6 := (uint16(c.G) * 3) / 255
	b6 := (uint16(c.B) * 3) / 255
	i := 36*r6 + 6*g6 + b6
	return []byte(fmt.Sprintf("\x1b[38;5;%dm", i))
}

func (c *color256) reset() []byte {
	return []byte("\x1b[0m")
}

func New256Colorer() colorer {
	return &color256{}
}
