package pixel

import (
	"time"
	"log"
	"fmt"
	"math"
	"sync"

	"github.com/tarm/serial"
)

type Pixel struct{
	sync.Mutex
	PortName      string
	PortSpeed     int
	lastPixelData PixelData
	serial        *serial.Port
}

type PixelData struct {
	Value      int
	Message    string
	Blink      int
	Brightness int
}

func (p *Pixel) Connect(){
	p.lastPixelData = PixelData{ -1, "", 0, 100 }

	c := &serial.Config{Name: p.PortName, Baud: p.PortSpeed}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Could not open port %s, %s", c.Name, err)
	}
	p.serial = s

	// port not opened before 1500 milliseconds pause
	time.Sleep(1500 * time.Millisecond)
}

func (p Pixel) GetStatus() PixelData{
	return p.lastPixelData
}

func (p *Pixel) SetStatus(pd PixelData){
	p.Lock()
	animationDuration := 3000 // ms

	switch pd.Blink {
	case 1:
		pd.Value += 100
	case 2:
		pd.Value += 200
	}

	delta := pd.Value - p.lastPixelData.Value
	stepTime := float64(animationDuration) / math.Abs(float64(delta))

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
				p.sendSerial(PixelData{i, "", 0, pd.Brightness })
				time.Sleep(time.Millisecond * time.Duration(stepTime))
			}
		}
	} else {
		// sharp switch color
		if(pd.Blink == 0 && p.lastPixelData.Value > 0 && pd.Value > 0) {
			for i := 0; i < 3; i++ {
				p.sendSerial(PixelData{pd.Value, "", 0, pd.Brightness })
				time.Sleep(time.Millisecond * 250)
				p.sendSerial(PixelData{p.lastPixelData.Value, "", 0, pd.Brightness })
				time.Sleep(time.Millisecond * 250)
			}
			p.sendSerial(PixelData{pd.Value, "", 0, pd.Brightness })
		}
	}

	time.Sleep(100 * time.Duration(stepTime))
	log.Printf("Pixel.setStatus: %v\n", pd)
	p.sendSerial(pd)

	p.lastPixelData = pd

	// if success value, turn off led
	if pd.Value == 100{
		time.Sleep(time.Millisecond * 5000)
		p.sendSerial(PixelData{ -1, "", 0, 100 })
		p.lastPixelData = PixelData{ -1, "", 0, 100 }
	}

	time.Sleep(1000 * time.Millisecond)
	p.Unlock()
}

func (p Pixel) sendSerial (pd PixelData) (int, error){
	command := fmt.Sprintf("%d|%s|%d\n",pd.Value, pd.Message, pd.Brightness)
	//log.Println(command)
	n, err := p.serial.Write([]byte(command))
	if err != nil {
		log.Fatalf("Could not write to port, %s", err)
	}
	return n, err
}
