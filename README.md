[![Build Status](https://travis-ci.org/popstas/pixel-server.svg?branch=travis-release)](https://travis-ci.org/popstas/pixel-server)

Server for send signals to [arduino pixel_meter](https://github.com/popstas/arduino-pixel-meter).

# Usage
- Start server, it will default listen at *:8080
- Send POST request to `/status` like this:
```
http -f POST http://localhost:8080/status value=50 message='first string\\second string' blink=2
```
- Or configure your [Kapacitor](https://github.com/influxdata/kapacitor) to `/kapacitor` like this:
```
data
    |alert()
        .post('http://localhost:8080/kapacitor')
```

# Configure server

### Command-line parameters
```
pixel-server \
--web-host="" \
--web-port=8080 \
--serial-port=COM3 \
--serial-speed=9600
```

### Environment variables
```
PIXEL_SERVER_SERIAL_PORT=COM3 \
PIXEL_SERVER_SERIAL_SPEED=9600 \
PIXEL_SERVER_WEB_HOST= \
PIXEL_SERVER_WEB_PORT=8080 \
pixel-server
```

Command-line parameters has priority over environment variables.

### Request parameters for /status
- `value` - value of signal, required,  
   0 to 100 (red to green),  
   -1 for off led
- `message` - message for 16x2 display, lines should be splitted with \ symbol, default no message
- `blink` - blink state, default 0,  
   0 for not blinking,  
   1 for blink 3 times and back to previous state,  
   2 for persistent blinking
- `brightness` - led brightness, 0 to 100, default 100

## Behaviour
If status changes from red to green, will be used smooth color change.
If status changes from green to red, color will changed with blinking with last color.