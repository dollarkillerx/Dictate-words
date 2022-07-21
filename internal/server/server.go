package server

import (
	"github.com/gin-gonic/gin"

	"os"
	"sync"
	"time"
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

	go s.gc()

	s.app.GET("/", s.showIdx)
	s.app.POST("generate_tts", s.generateTTS)
	s.app.GET("download_tts/:tts_id", s.downloadTTS)
}

func (s *Server) gc() {
	for {
		time.Sleep(time.Minute * 50)
		s.gcCore()
	}
}

func (s *Server) gcCore() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, v := range s.cache {
		if v.Expiration < time.Now().Unix() {
			os.Remove(v.Filepath)
			delete(s.cache, k)
		}
	}
}
