package main

import (
	"github.com/dollarkillerx/Dictate-words/internal/server"

	"log"
)

func main() {
	server := server.NewServer()
	err := server.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
