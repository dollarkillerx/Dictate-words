package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"sync"
)

type Server struct {
	app *gin.Engine

	cache map[string]WordCache
	mu    sync.Mutex
}

func NewServer() *Server {
	ser := Server{
		app:   gin.Default(),
		cache: map[string]WordCache{},
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
	s.app.GET("/", s.showIdx)
	s.app.POST("generate_tts", s.generateTTS)
	s.app.GET("download_tts/:tts_id", s.downloadTTS)
}

//
//ffmpeg -i "concat:./stats/start.mp3|./stats/temporary/cbcafgj06pn1sp0edtbg_cbcafgj06pn1sp0edtd0.mp3|./stats/ting.mp3|./stats/temporary/cbcafgj06pn1sp0edtbg_cbcafgj06pn1sp0edtd0.mp3|./stats/ting.mp3|./stats/temporary/cbcafgj06pn1sp0edtbg_cbcafgj06pn1sp0edtd0.mp3|./stats/ting.mp3" -acodec copy cbcafgj06pn1sp0edtbg.mp3 -y
//ffmpeg -i "./stats/start.mp3|./stats/temporary/cbcafgj06pn1sp0edtbg_cbcafgj06pn1sp0edtd0.mp3" -acodec copy cbcafgj06pn1sp0edtbg.mp3 -y
