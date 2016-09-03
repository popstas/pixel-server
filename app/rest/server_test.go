package rest

import (
	"testing"
	"net/http/httptest"
	"os"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/popstas/pixel-server/app/pixel"
	"net/url"
	"strings"
)

func getServer() Server{
	p := pixel.SerialPixel{Serial: pixel.NullWriter{}, Testing: true}
	s := Server{Pixel: p}
	return s
}

func sendKapacitorRequestFromFile(fileName string, r *gin.Engine) (*httptest.ResponseRecorder){
	response := httptest.NewRecorder()
	f, _ := os.Open(fileName)
	request := httptest.NewRequest("POST", "/kapacitor", f)
	r.ServeHTTP(response, request)
	return response
}

func TestKapacitor(t *testing.T) {
	var response *httptest.ResponseRecorder
	s := getServer()
	r := s.GetEngine()
	//httptest.NewServer(e)

	response = sendKapacitorRequestFromFile("fixtures/kapacitor_warning.json", r)
	log.Printf("Response: %v", response)

	response = sendKapacitorRequestFromFile("fixtures/kapacitor_critical.json", r)
	log.Printf("Response: %v", response)

	response = sendKapacitorRequestFromFile("fixtures/kapacitor_ok.json", r)
	log.Printf("Response: %v", response)
}

func TestStatus(t *testing.T) {
	response := httptest.NewRecorder()
	s := getServer()
	r := s.GetEngine()

	form := url.Values{}
	form.Add("value", "100")
	form.Add("message", "text\\second line")
	form.Add("blink", "0")
	form.Add("brightness", "100")

	request := httptest.NewRequest("POST", "/status", strings.NewReader(form.Encode()))
	request.PostForm = form
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(response, request)

	form = url.Values{}
	form.Add("value", "100")
	form.Add("message", "text\\second line")
	form.Add("blink", "0")

	request = httptest.NewRequest("POST", "/status", strings.NewReader(form.Encode()))
	request.PostForm = form
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(response, request)
}