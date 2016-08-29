package main

import (
	"time"
	"fmt"
)

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
