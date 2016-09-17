package outputs

import (
	"log"
	"time"
	"sync"
	"github.com/justincampbell/anybar"
	"github.com/mitchellh/go-ps"
	"errors"
)

type AnyBar struct {
	sync.Mutex
	Port          int
	lastPixelData PixelData
	Testing       bool
}

func checkAnyBar() bool {
	proc, _ := ps.Processes()
	for _, p := range proc {
		if p.Executable() == "AnyBar" {
			return true
		}
	}
	return false
}

func CreateAnyBar(port int) (error, *AnyBar) {
	if !checkAnyBar() {
		return errors.New("AnyBar not running"), nil
	}
	return nil, &AnyBar{Port: port, lastPixelData: PixelData{-1, "", 0, 100 }}
}

func (p *AnyBar) GetState() PixelData {
	return p.lastPixelData
}

func (p *AnyBar) SetState(pd PixelData) {
	p.Lock()
	log.Printf("[AnyBar] Set state: %v\n", pd)

	cmds := p.BuildSetStateCommands(pd)
	p.ExecCommands(cmds)
	p.lastPixelData = pd

	// if success value, turn off led
	if pd.Value == 100{
		p.lastPixelData = PixelData{ -1, "", 0, 100 }
	}

	p.Unlock()
}

func (p AnyBar) BuildSetStateCommands (pd PixelData) ([]StateCommand){
	cmds := []StateCommand{}

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

func (p AnyBar) ExecCommands(cmds []StateCommand){
	for _, c := range cmds{
		if c.PauseBefore > 0 && !p.Testing{
			time.Sleep(time.Millisecond * time.Duration(c.PauseBefore))
		}

		p.sendAnyBar(c.Pd)

		if c.PauseAfter > 0 && !p.Testing{
			time.Sleep(time.Millisecond * time.Duration(c.PauseAfter))
		}
	}
}

func (p AnyBar) sendAnyBar(pd PixelData) {
	if pd.Value > 100{
		pd.Value -= 100
	}
	if pd.Value > 200{
		pd.Value -= 200
	}

	switch {
	case pd.Value >= 1 && pd.Value < 50:
		anybar.Red()
	case pd.Value >= 50 && pd.Value < 100:
		anybar.Yellow()
	case pd.Value == 100:
		anybar.Green()
	case pd.Value == -1:
		anybar.White()
	}
}