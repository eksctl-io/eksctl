package lol

import (
	"fmt"
	"math"
)

type truecolor struct {
	R, G, B uint8
}

func (c *truecolor) rainbow(freq float64, i float64) {
	c.R = uint8(math.Floor(math.Sin(freq*i+0)*127)) + 128
	c.G = uint8(math.Floor(math.Sin(freq*i+2.0*math.Pi/3.0)*127)) + 128
	c.B = uint8(math.Floor(math.Sin(freq*i+4.0*math.Pi/3.0)*127)) + 128
}

func (c *truecolor) format() []byte {
	return []byte(fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B))
}

func (c *truecolor) reset() []byte {
	return []byte("\x1b[0m")
}

func newTruecolorColorer() colorer {
	return &truecolor{}
}
