package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/popstas/pixel-server/app/rest"
	"github.com/popstas/pixel-server/app/outputs"
	"log"
	"runtime"
)

var opts struct {
	SerialPort  string `long:"serial-port" env:"PIXEL_SERVER_SERIAL_PORT" description:"serial port name or path" default:"COM3"`
	SerialSpeed int    `long:"serial-speed" env:"PIXEL_SERVER_SERIAL_SPEED" description:"serial port speed" default:"9600"`
	AnyBarPort  int    `long:"anybar-port" env:"PIXEL_SERVER_ANYBAR_PORT" description:"anybar port" default:"1738"`
	WebHost     string `long:"web-host" env:"PIXEL_SERVER_WEB_HOST" description:"hostname for bind server" default:""`
	WebPort     int    `long:"web-port" env:"PIXEL_SERVER_WEB_PORT" description:"port for bind server" default:"8246"`
	Brightness  int    `long:"brightness" env:"PIXEL_SERVER_BRIGHTNESS" description:"default pixel's brightness" default:"100"`
}

func main() {

	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	pixels := []outputs.Pixel{}

	// init SerialPixel
	if err, p := outputs.CreateSerialPixel(opts.SerialPort, opts.SerialSpeed); err != nil {
		log.Println(err)
	} else {
		pixels = append(pixels, p)
	}

	// init AnyBar
	if runtime.GOOS == "darwin" {
		if err, anybar := outputs.CreateAnyBar(opts.AnyBarPort); err != nil {
			log.Println(err)
		} else {
			pixels = append(pixels, anybar)
		}
	}

	// init Web Server
	server := rest.Server{
		HostPort: fmt.Sprintf("%s:%d", opts.WebHost, opts.WebPort),
		Pixels: pixels,
		DefaultBrightness: opts.Brightness,
	}

	server.Run()
}
