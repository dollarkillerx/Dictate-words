package test

import (
	"bytes"
	"github.com/dollarkillerx/Dictate-words/pkg/models"
	"github.com/dollarkillerx/urllib"
	"github.com/viert/go-lame"
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

	var resp models.TtsResp
	err = urllib.Post("http://tts_api.mechat.live/google_tts").SetJsonObject(models.TtsModel{
		Text: "こんにちは、世界。",
		Lang: "ja",
	}).FromJsonByCode(&resp, 200)
	if err != nil {
		log.Fatalln(err)
	}

	code, bt, err := urllib.Get(resp.Url).ByteOriginal()
	if err != nil {
		log.Fatalln(err)
	}

	if code != 200 {
		log.Fatalln(string(bt))
	}

	create, err := os.Create("xxx.mp3")
	if err != nil {
		log.Fatalln(err)
	}
	defer create.Close()

	enc := lame.NewEncoder(create)
	defer enc.Close()

	r := bytes.NewReader(start)
	r.WriteTo(enc)

	for i := 0; i < 3; i++ {
		r := bytes.NewReader(bt)
		r.WriteTo(enc)

		r = bytes.NewReader(ting)
		r.WriteTo(enc)
	}

}
