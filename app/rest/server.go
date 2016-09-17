package rest

import (
	"time"
	"log"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/popstas/pixel-server/app/outputs"
	"github.com/popstas/pixel-server/app/kapacitor"
)

type Server struct {
	HostPort string
	Pixels    []outputs.Pixel
}

func (s Server) GetEngine() *gin.Engine {
	router := gin.Default()

	router.POST("/status", s.statusHandler)
	router.POST("/kapacitor", s.kapacitorHandler)

	return router
}

func (s *Server) Run() {
	gin.SetMode(gin.ReleaseMode)
	router := s.GetEngine()

	s.setState(outputs.PixelData{ 100, fmt.Sprintf("server started\\%s", s.HostPort), 1, 20 })
	time.Sleep(2000 * time.Millisecond)
	s.setState(outputs.PixelData{ -1, "", 0, 100 })

	log.Fatal(router.Run(s.HostPort))
}

func (s Server) setState (pd outputs.PixelData){
	for _, pixel := range s.Pixels {
		go pixel.SetState(pd)
	}
}

func (s *Server) statusHandler(c *gin.Context) {
	var pd outputs.PixelData
	var err error

	pd.Value, _ = strconv.Atoi(c.PostForm("value"))
	pd.Message = c.PostForm("message")
	pd.Blink, _ = strconv.Atoi(c.PostForm("blink"))
	pd.Brightness, err = strconv.Atoi(c.PostForm("brightness"))
	if err != nil{
		pd.Brightness = 100
	}

	s.setState(pd)
}

func (s *Server) kapacitorHandler(c *gin.Context) {
	ad := kapacitor.KapacitorAlertData{}
	err := c.BindJSON(&ad)
	if err != nil {
		log.Fatalf("Could not decode KapacitorAlertData, %s", err)
	}

	var pd outputs.PixelData
	pd.Brightness = 100

	switch ad.Level {
	case kapacitor.OKAlert:
		pd.Value = 100
		//pd.Blink = 1
	case kapacitor.InfoAlert:
		pd.Value = -1
	case kapacitor.WarnAlert:
		pd.Value = 50
	case kapacitor.CritAlert:
		pd.Value = 1
		//pd.Blink = 2
	}

	data := ad.Data.Series[0]
	pd.Message = fmt.Sprintf("%s\\%s: %v", data.Tags.Host, data.Name, data.Values[0][1]) // data.Values[1]

	s.setState(pd)
}
