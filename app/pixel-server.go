package main

import (
	"time"
	"net/http"
	"log"
	"fmt"
	"strconv"
	"encoding/json"
	"math"

	"github.com/tarm/serial"
)

type PixelServer struct {
	Config        Config
	Serial        *serial.Port
	LastPixelData PixelData
}

type PixelData struct {
	Value      int
	Message    string
	Blink      int
	Brightness int
}

func (ps PixelServer) statusHandler(w http.ResponseWriter, r *http.Request) {
	if(r.Method != "POST") {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()

	var pd PixelData

	pd.Value, _ = strconv.Atoi(r.FormValue("value"))
	pd.Message = r.FormValue("message")
	pd.Blink, _ = strconv.Atoi(r.FormValue("blink"))
	pd.Brightness, _ = strconv.Atoi(r.FormValue("brightness"))

	pixelServer.setStatus(pd)
}

func (ps PixelServer) kapacitorHandler(w http.ResponseWriter, r *http.Request) {
	if(r.Method != "POST") {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	ad := KapacitorAlertData{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ad)
	if err != nil {
		log.Fatalf("Could not decode kapacitor AlertData, %s", err)
	}

	var pd PixelData
	pd.Brightness = 100

	switch ad.Level {
	case OKAlert:
		pd.Value = 100
		//pd.Blink = 1
	case InfoAlert:
		pd.Value = -1
	case WarnAlert:
		pd.Value = 50
	case CritAlert:
		pd.Value = 1
		pd.Blink = 2
	}

	data := ad.Data.Series[0]
	pd.Message = fmt.Sprintf("%s\\%s: %v", data.Tags.Host, data.Name, data.Values[0][1]) // data.Values[1]

	pixelServer.setStatus(pd)
}

func (ps *PixelServer) setStatus (pd PixelData){
	animationDuration := 1000 // ms

	switch pd.Blink {
	case 1:
		pd.Value += 100
	case 2:
		pd.Value += 200
	}

	delta := pd.Value - ps.LastPixelData.Value
	stepTime := float64(animationDuration) / math.Abs(float64(delta))

	var step int
	if delta > 0{
		step = 1
	} else {
		step = -1
	}

	if(pd.Blink == 0 && ps.LastPixelData.Value > 0 && pd.Value > 0) {
		for i := ps.LastPixelData.Value; i != pd.Value; i += step {
			ps.sendSerial(PixelData{i, "", 0, pd.Brightness })
			time.Sleep(time.Millisecond * time.Duration(stepTime))
		}
	}

	time.Sleep(100 * time.Duration(stepTime))
	log.Printf("setStatus: %v\n", pd)
	ps.sendSerial(pd)

	// if success value, turn off led
	if pd.Value == 100{
		time.Sleep(time.Millisecond * 5000)
		ps.sendSerial(PixelData{ -1, "", 0, 100 })
	}

	time.Sleep(1000 * time.Millisecond)
	ps.LastPixelData = pd
}

func (ps PixelServer) sendSerial (pd PixelData) (int, error){
	command := fmt.Sprintf("%d|%s|%d\n",pd.Value, pd.Message, pd.Brightness)
	//fmt.Println(command)
	n, err := ps.Serial.Write([]byte(command))
	if err != nil {
		log.Fatalf("Could not write to port, %s", err)
	}
	return n, err
}
