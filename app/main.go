package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/popstas/pixel-server/app/rest"
	"github.com/popstas/pixel-server/app/pixel"
)

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

	p := pixel.SerialPixel{PortName: opts.SerialPort, PortSpeed: opts.SerialSpeed}
	p.Connect()

	server := rest.Server{
		HostPort: fmt.Sprintf("%s:%d", opts.WebHost, opts.WebPort),
		Pixel: p,
	}
	server.Run()
}
