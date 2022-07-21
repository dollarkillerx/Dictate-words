package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

type Server struct {
	app *gin.Engine
}

func NewServer() *Server {
	ser := Server{
		app: gin.Default(),
	}

	return &ser
}

func (s *Server) Run() error {
	s.router()
	return s.app.Run("0.0.0.0:8474")
}

func (s *Server) router() {
	s.app.Use(Cors())

	os.Setenv("HTTPS_PROXY", "http://192.168.31.192:8081")
	os.Setenv("HTTP_PROXY", "http://192.168.31.192:8081")

	log.Println("===========")
	log.Println(os.Getenv("HTTPS_PROXY"))
	s.app.POST("generate_tts", s.generateTTS)
}

//
//ffmpeg -i "concat:./stats/start.mp3|./stats/temporary/cbcafgj06pn1sp0edtbg_cbcafgj06pn1sp0edtd0.mp3|./stats/ting.mp3|./stats/temporary/cbcafgj06pn1sp0edtbg_cbcafgj06pn1sp0edtd0.mp3|./stats/ting.mp3|./stats/temporary/cbcafgj06pn1sp0edtbg_cbcafgj06pn1sp0edtd0.mp3|./stats/ting.mp3" -acodec copy cbcafgj06pn1sp0edtbg.mp3 -y
//ffmpeg -i "./stats/start.mp3|./stats/temporary/cbcafgj06pn1sp0edtbg_cbcafgj06pn1sp0edtd0.mp3" -acodec copy cbcafgj06pn1sp0edtbg.mp3 -y
