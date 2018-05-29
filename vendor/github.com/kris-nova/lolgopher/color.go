package lol

// colorer is an interface that can be used to colorize an io.Writer. Each
// colorer is designed to support a specific terminal color schema (256, 0, etc).
type colorer interface {

	// Rainbow will calculate the current color spectrum for RGB colors only.
	// Each implementation other than truecolor will need to implement rainbow
	// and handle it accordingly.
	rainbow(freq float64, i float64)

	// Return the ASCII escape code for this color for this particular byte
	format() []byte

	// Reset the most recent ASCII code opened
	reset() []byte
}
