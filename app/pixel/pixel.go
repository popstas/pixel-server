package pixel

import (
	"time"
	"log"
	"fmt"
	"math"
	"sync"

	"github.com/tarm/serial"
	"io"
)

type Pixel interface {
	SetState(pd PixelData)
	GetState() (pd PixelData)
}

type SerialPixel struct{
	sync.Mutex
	PortName      string
	PortSpeed     int
	lastPixelData PixelData
	Serial        io.Writer
	Testing       bool
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

func (p *SerialPixel) Connect(){
	p.lastPixelData = PixelData{ -1, "", 0, 100 }

	c := &serial.Config{Name: p.PortName, Baud: p.PortSpeed}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Could not open port %s, %s", c.Name, err)
	}
	p.Serial = s

	// port not opened before 1500 milliseconds pause
	time.Sleep(1500 * time.Millisecond)
}

func (p SerialPixel) GetState() PixelData{
	return p.lastPixelData
}

func (p *SerialPixel) SetState(pd PixelData){
	p.Lock()

	switch pd.Blink {
	case 1:
		pd.Value += 100
	case 2:
		pd.Value += 200
	}

	cmds := p.BuildSetStateCommands(pd)
	if !p.Testing{
		p.ExecCommands(cmds)
	}

	log.Printf("Pixel.SetState: %v\n", pd)

	p.lastPixelData = pd

	// if success value, turn off led
	if pd.Value == 100{
		p.lastPixelData = PixelData{ -1, "", 0, 100 }
	}

	p.Unlock()
}

func (p SerialPixel) BuildSetStateCommands (pd PixelData) ([]StateCommand){
	cmds := []StateCommand{}
	animationDuration := 3000 // ms

	delta := pd.Value - p.lastPixelData.Value
	stepTime := int(float64(animationDuration) / math.Abs(float64(delta)))

	var step int
	if delta > 0{
		step = 1
	} else {
		step = -1
	}

	if delta > 0{
		// smooth switch color
		if(pd.Blink == 0 && p.lastPixelData.Value > 0 && pd.Value > 0) {
			for i := p.lastPixelData.Value; i != pd.Value; i += step {
				cmds = append(cmds, StateCommand{
					Pd: PixelData{i, "", 0, pd.Brightness },
					PauseAfter: stepTime,
				})
			}
		}
	} else {
		// sharp switch color
		if(pd.Blink == 0 && p.lastPixelData.Value > 0 && pd.Value > 0) {
			for i := 0; i < 3; i++ {
				cmds = append(cmds, StateCommand{
					Pd: PixelData{pd.Value, "", 0, pd.Brightness },
					PauseAfter: 250,
				})
				cmds = append(cmds, StateCommand{
					Pd: PixelData{p.lastPixelData.Value, "", 0, pd.Brightness },
					PauseAfter: 250,
				})
			}
			cmds = append(cmds, StateCommand{
				Pd: PixelData{pd.Value, "", 0, pd.Brightness },
				PauseAfter: 250,
			})
		}
	}

	cmds = append(cmds, StateCommand{
		Pd: pd,
		PauseAfter: 1000,
	})

	// if success value, turn off led
	if pd.Value == 100 {
		cmds = append(cmds, StateCommand{
			Pd: PixelData{-1, "", 0, 100 },
			PauseBefore: 5000,
		})
	}

	return cmds
}

func (p SerialPixel) ExecCommands(cmds []StateCommand){
	for _, c := range cmds{
		if c.PauseBefore > 0{
			time.Sleep(time.Millisecond * time.Duration(c.PauseBefore))
		}

		p.sendSerial(c.Pd)

		if c.PauseAfter > 0{
			time.Sleep(time.Millisecond * time.Duration(c.PauseAfter))
		}
	}
}

func (p SerialPixel) sendSerial (pd PixelData) (int, error){
	command := fmt.Sprintf("%d|%s|%d\n",pd.Value, pd.Message, pd.Brightness)
	//log.Println(command)
	n, err := p.Serial.Write([]byte(command))
	if err != nil {
		log.Fatalf("Could not write to port, %s", err)
	}
	return n, err
}
