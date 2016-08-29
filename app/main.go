package main

import (
	"log"
	"net/http"
	"time"
	"fmt"
	"flag"

	"github.com/tarm/serial"
)

var pixelServer PixelServer

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
