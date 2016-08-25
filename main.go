package main

import (
	"log"
	"net/http"

	"github.com/tarm/serial"
	"time"
	"fmt"
	"strconv"
	"encoding/json"
)

const (
	SerialPort  = "COM3"
	SerialSpeed = 9600
	WebHost     = ""
	WebPort     = 8080
)

type PixelServer struct {

	Serial *serial.Port
}

type PixelData struct {
	Value int
	Message string
	Blink int
	Brightness int
}

// from https://github.com/influxdata/kapacitor/blob/master/alert.go#L87
type KapacitorAlertData struct {
	ID       string              `json:"id"`
	Message  string              `json:"message"`
	Details  string              `json:"details"`
	Time     time.Time           `json:"time"`
	Duration time.Duration       `json:"duration"`
	Level    KapacitorAlertLevel `json:"level"`
	Data struct {
		Series []struct {
			Columns []string `json:"columns"`
			Name    string   `json:"name"`
			Tags    struct {
				Host string `json:"host"`
			} `json:"tags"`
			Values [][]interface{} `json:"values"`
		} `json:"series"`
	} `json:"data"`
}

type KapacitorAlertLevel int

const (
	OKAlert KapacitorAlertLevel = iota
	InfoAlert
	WarnAlert
	CritAlert
)

func (l KapacitorAlertLevel) String() string {
	switch l {
	case OKAlert:
		return "OK"
	case InfoAlert:
		return "INFO"
	case WarnAlert:
		return "WARNING"
	case CritAlert:
		return "CRITICAL"
	default:
		panic("unknown AlertLevel")
	}
}

func (l KapacitorAlertLevel) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *KapacitorAlertLevel) UnmarshalText(text []byte) error {
	s := string(text)
	switch s {
	case "OK":
		*l = OKAlert
	case "INFO":
		*l = InfoAlert
	case "WARNING":
		*l = WarnAlert
	case "CRITICAL":
		*l = CritAlert
	default:
		return fmt.Errorf("unknown AlertLevel %s", s)
	}
	return nil
}

var pixelServer PixelServer

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
	pd.Brightness = 100

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
		pd.Blink = 1
	case InfoAlert:
		pd.Value = -1
	case WarnAlert:
		pd.Value = 50
	case CritAlert:
		pd.Value = 1
		pd.Blink = 2
	}

	data := ad.Data.Series[0]
	fmt.Printf("%v\n", data)
	fmt.Printf("%v\n", data.Values)
	fmt.Printf("%v\n", data.Values[0][1])
	pd.Message = fmt.Sprintf("%s\\%s: %v", data.Tags.Host, data.Name, data.Values[0][1]) // data.Values[1]

	pixelServer.setStatus(pd)
}

func (ps PixelServer) setStatus (pd PixelData) (int, error){
	switch pd.Blink {
	case 1:
		pd.Value += 100
	case 2:
		pd.Value += 200
	}

	command := fmt.Sprintf("%d|%s|%d\n",pd.Value, pd.Message, pd.Brightness)
	n, err := ps.Serial.Write([]byte(command))
	if err != nil {
		log.Fatalf("Could not write to port, %s", err)
	}

	return n, err
}

func main() {
	c := &serial.Config{Name: SerialPort, Baud: SerialSpeed}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Could not open port %s, %s", c.Name, err)
	}
	pixelServer.Serial = s
	hostPort := fmt.Sprintf("%s:%d", WebHost, WebPort)

	// port not opened before 1500 milliseconds pause
	time.Sleep(1500 * time.Millisecond)

	initCommand := fmt.Sprintf("200|server started\\%s|50\n", hostPort)
	n, err := s.Write([]byte(initCommand))
	if err != nil {
		log.Fatal(err)
	}
	_ = n

	time.Sleep(2000 * time.Millisecond)
	s.Write([]byte("-1\n"))

	http.HandleFunc("/status", pixelServer.statusHandler)
	http.HandleFunc("/kapacitor", pixelServer.kapacitorHandler)
	http.ListenAndServe(hostPort, nil)
}
