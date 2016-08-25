package main

import (
	"log"
	"net/http"

	"github.com/tarm/serial"
	"time"
	"fmt"
	"strconv"
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

var pixelServer PixelServer

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if(r.Method == "POST"){
		r.ParseForm()

		value, _ := strconv.Atoi(r.FormValue("value"))
		message := r.FormValue("message")
		blink, _ := strconv.Atoi(r.FormValue("blink"))
		brightness := 100

		switch blink {
		case 1:
			value += 100
		case 2:
			value += 200
		}

		command := fmt.Sprintf("%d|%s|%d\n",value, message, brightness)
		n, err := pixelServer.Serial.Write([]byte(command))
		if err != nil {
			log.Fatal(err)
		}
		_ = n
	} else{
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
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

	http.HandleFunc("/status", statusHandler)
	http.ListenAndServe(hostPort, nil)
}
