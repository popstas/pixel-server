package main

import (
	"log"
	"net/http"
	"time"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/tarm/serial"
)

var pixelServer PixelServer

var opts struct {
	SerialPort  string `long:"serial-port" env:"PIXEL_SERVER_SERIAL_PORT" description:"serial port name or path" default:"COM3"`
	SerialSpeed int    `long:"serial-speed" env:"PIXEL_SERVER_SERIAL_SPEED" description:"serial port speed" default:"9600"`
	WebHost     string `long:"web-host" env:"PIXEL_SERVER_WEB_HOST" description:"hostname for bind server" default:""`
	WebPort     int    `long:"web-port" env:"PIXEL_SERVER_WEB_PORT" description:"port for bind server" default:"8080"`
}

func main() {

	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	c := &serial.Config{Name: opts.SerialPort, Baud: opts.SerialSpeed}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Could not open port %s, %s", c.Name, err)
	}
	pixelServer.Serial = s
	hostPort := fmt.Sprintf("%s:%d", opts.WebHost, opts.WebPort)

	// port not opened before 1500 milliseconds pause
	time.Sleep(1500 * time.Millisecond)

	pixelServer.setStatus(PixelData{ 100, fmt.Sprintf("server started\\%s", hostPort), 1, 20 })
	time.Sleep(2000 * time.Millisecond)
	pixelServer.setStatus(PixelData{ -1, "", 0, 100 })

	http.HandleFunc("/status", pixelServer.statusHandler)
	http.HandleFunc("/kapacitor", pixelServer.kapacitorHandler)
	http.ListenAndServe(hostPort, nil)
}
