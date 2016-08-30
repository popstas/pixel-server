package rest

import (
	"time"
	"net/http"
	"log"
	"fmt"
	"strconv"
	"encoding/json"

	"github.com/popstas/pixel-server/app/pixel"
	"github.com/popstas/pixel-server/app/kapacitor"
)

type Server struct {
	Pixel         pixel.Pixel
	HostPort      string
}

func (ps Server) Run() {
	ps.Pixel.SetStatus(pixel.PixelData{ 100, fmt.Sprintf("server started\\%s", ps.HostPort), 1, 20 })
	time.Sleep(2000 * time.Millisecond)
	ps.Pixel.SetStatus(pixel.PixelData{ -1, "", 0, 100 })

	http.HandleFunc("/status", ps.statusHandler)
	http.HandleFunc("/kapacitor", ps.kapacitorHandler)
	http.ListenAndServe(ps.HostPort, nil)
}

func (ps *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	if(r.Method != "POST") {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()

	var pd pixel.PixelData
	var err error

	pd.Value, _ = strconv.Atoi(r.FormValue("value"))
	pd.Message = r.FormValue("message")
	pd.Blink, _ = strconv.Atoi(r.FormValue("blink"))
	pd.Brightness, err = strconv.Atoi(r.FormValue("brightness"))
	if err != nil{
		pd.Brightness = 100
	}

	go ps.Pixel.SetStatus(pd)
}

func (ps Server) kapacitorHandler(w http.ResponseWriter, r *http.Request) {
	if(r.Method != "POST") {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	ad := kapacitor.KapacitorAlertData{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ad)
	if err != nil {
		log.Fatalf("Could not decode kapacitor AlertData, %s", err)
	}

	var pd pixel.PixelData
	pd.Brightness = 100

	switch ad.Level {
	case kapacitor.OKAlert:
		pd.Value = 100
		//pd.Blink = 1
	case kapacitor.InfoAlert:
		pd.Value = -1
	case kapacitor.WarnAlert:
		pd.Value = 50
	case kapacitor.CritAlert:
		pd.Value = 1
		//pd.Blink = 2
	}

	data := ad.Data.Series[0]
	pd.Message = fmt.Sprintf("%s\\%s: %v", data.Tags.Host, data.Name, data.Values[0][1]) // data.Values[1]

	go ps.Pixel.SetStatus(pd)
}

