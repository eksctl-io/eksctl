package lol

type color0 struct {
	R, G, B uint8
}

func (c *color0) rainbow(freq float64, i float64) {
	// No calculation for no color support
}

func (c *color0) format() []byte {
	return []byte("")
}

func (c *color0) reset() []byte {
	return []byte("")
}

func New0Colorer() colorer {
	return &color0{}
}
