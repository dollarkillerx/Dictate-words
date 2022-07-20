package test

import (
	"bytes"
	"github.com/dollarkillerx/Dictate-words/pkg/models"
	"github.com/dollarkillerx/urllib"
	"github.com/hyacinthus/mp3join"
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

	joiner := mp3join.New()

	err = joiner.Append(bytes.NewBuffer(start))
	if err != nil {
		log.Fatalln(err)
	}

	// readers is the input mp3 files
	for i := 0; i < 3; i++ {
		err = joiner.Append(bytes.NewBuffer(bt))
		if err != nil {
			log.Fatalln(err)
		}
		err = joiner.Append(bytes.NewBuffer(ting))
		if err != nil {
			log.Fatalln(err)
		}
	}

	dest := joiner.Reader()
	dest.WriteTo(create)
}
