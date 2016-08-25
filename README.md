Server for send signals to [arduino pixel_meter](https://github.com/popstas/arduino-pixel-meter).

# Usage
- Start server, it will listen at *:8080
- Send signal like this:
```
http -f POST http://localhost:8080/status value=50 message='first string\\second string' blink=2
```

Requests should be POST to /status

### Request parameters
- `value` - value of signal, 0 to 100 (red to green), -1 for off led
- `message` - message for 16x2 display, lines should be splitted with \ symbol
- `blink` - blink state, 0 for not blinking, 1 for blink 3 times and back to previous state, 2 for persistent blinking
