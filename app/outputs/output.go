package outputs

import "log"

type Pixel interface {
	SetState(pd PixelData)
	GetState() (pd PixelData)
}

type PixelData struct {
	Value      int
	Message    string
	Blink      int
	Brightness int
}

type StateCommand struct {
	Pd          PixelData
	PauseBefore int
	PauseAfter  int
}

type NullWriter struct {
	LastSerialMessage []byte
}

func (w NullWriter) Write (b []byte) (int, error){
	w.LastSerialMessage = b
	log.Printf("Write bytes: %s", b)
	return 0, nil
}

