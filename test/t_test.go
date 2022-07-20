package test

import (
	"fmt"
	"github.com/dollarkillerx/Dictate-words/utils/googletts"
	"github.com/dollarkillerx/urllib"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestM2(t *testing.T) {
	start, err := ioutil.ReadFile("../stats/start.mp3")
	if err != nil {
		log.Fatalln(err)
	}

	ting, err := ioutil.ReadFile("../stats/ting.mp3")
	if err != nil {
		log.Fatalln(err)
	}

	url, err := googletts.GetTTSURL("こんにちは、世界。", "ja")
	if err != nil {
		panic(err)
	}
	fmt.Println(url) // => https://translate.google.com/translate_tts?client=t&ie=UTF-8&q=Hello%2C+world.&textlen=13&tk=368668.249914&tl=en

	code, bytes, err := urllib.Get(url).ByteOriginal()
	if err != nil {
		log.Fatalln(err)
	}
	if code != 200 {
		log.Fatalln(string(bytes))
	}

	create, err := os.Create("xxx.mp3")
	if err != nil {
		log.Fatalln(err)
	}
	defer create.Close()

	create.Write(start)

	for i := 0; i < 3; i++ {
		create.Write(bytes)
		create.Write(ting)
	}

}
