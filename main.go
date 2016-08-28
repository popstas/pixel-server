package main

import (
	"log"
	"net/http"

	"github.com/tarm/serial"
	"time"
	"fmt"
	"strconv"
	"encoding/json"
	"math"
	"flag"
	"os"
)

type Config struct {
	SerialPort string
	SerialSpeed int
	WebHost string
	WebPort int
}

type PixelServer struct {
	Config Config
	Serial *serial.Port
	LastPixelData PixelData
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

func init() {
	setStringEnvvar(&pixelServer.Config.SerialPort, "PIXEL_SERVER_SERIAL_PORT")
	setIntEnvvar(&pixelServer.Config.SerialSpeed, "PIXEL_SERVER_SERIAL_SPEED")
	setStringEnvvar(&pixelServer.Config.WebHost, "PIXEL_SERVER_WEB_HOST")
	setIntEnvvar(&pixelServer.Config.WebPort, "PIXEL_SERVER_WEB_PORT")

	flag.StringVar(&pixelServer.Config.SerialPort, "serial-port", "COM3", "serial port name or path")
	flag.IntVar(&pixelServer.Config.SerialSpeed, "serial-speed",  9600, "serial port speed")
	flag.StringVar(&pixelServer.Config.WebHost, "web-host",  "", "hostname for bind server")
	flag.IntVar(&pixelServer.Config.WebPort, "web-port",  8080, "port for bind server")
	flag.Parse()
}

func setIntEnvvar(v *int, envName string){
	envValue := os.Getenv(envName)
	if envValue != ""{
		if envIntValue, err := strconv.Atoi(envValue); err != nil{
			log.Fatalf("Cannot convert value of envvar %s to int: %s", envName, envValue)
		} else {
			*v = envIntValue
		}
	}
}

func setStringEnvvar(v *string, envName string){
	envValue := os.Getenv(envName)
	if envValue != ""{
		*v = envValue
	}
}

func main() {
	c := &serial.Config{Name: pixelServer.Config.SerialPort, Baud: pixelServer.Config.SerialSpeed}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Could not open port %s, %s", c.Name, err)
	}
	pixelServer.Serial = s
	hostPort := fmt.Sprintf("%s:%d", pixelServer.Config.WebHost, pixelServer.Config.WebPort)

	// port not opened before 1500 milliseconds pause
	time.Sleep(1500 * time.Millisecond)

	pixelServer.setStatus(PixelData{ 100, fmt.Sprintf("server started\\%s", hostPort), 1, 20 })
	time.Sleep(2000 * time.Millisecond)
	pixelServer.setStatus(PixelData{ -1, "", 0, 100 })

	http.HandleFunc("/status", pixelServer.statusHandler)
	http.HandleFunc("/kapacitor", pixelServer.kapacitorHandler)
	http.ListenAndServe(hostPort, nil)
}
