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

	log.Println("===========")
	log.Println(os.Getenv("HTTPS_PROXY"))
	s.app.POST("generate_tts", s.generateTTS)
}
