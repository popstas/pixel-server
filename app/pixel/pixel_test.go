package pixel

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func assertCommandCountAfterSetState(t *testing.T, p *SerialPixel, pd PixelData, expected int){
	cmds := p.BuildSetStateCommands(pd)
	p.SetState(pd)
	assert.Equal(t, expected, len(cmds))
}

func TestSetState(t *testing.T) {
	p := SerialPixel{Serial: NullWriter{}, Testing: true}

	// expect sendSerial called 1 time
	assertCommandCountAfterSetState(t, &p, PixelData{1, "", 0, 100}, 1)

	// expect sendSerial called 50 times (smooth)
	assertCommandCountAfterSetState(t, &p, PixelData{50, "", 0, 100}, 50)

	// expect sendSerial called 50 times then called sendSerial -1
	assertCommandCountAfterSetState(t, &p, PixelData{100, "", 0, 100}, 52)

	// expect sendSerial called 1 time (after off)
	assertCommandCountAfterSetState(t, &p, PixelData{50, "", 0, 100}, 1)

	// expect sendSerial called 8 times (blinking)
	assertCommandCountAfterSetState(t, &p, PixelData{1, "", 0, 100}, 8)
}