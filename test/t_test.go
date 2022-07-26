package test

import (
	"github.com/dollarkillerx/Dictate-words/pkg/models"
	"github.com/dollarkillerx/urllib"
	"log"
	"os"
	"testing"
)

func TestM2(t *testing.T) {
	//start, err := ioutil.ReadFile("../stats/start.mp3")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//ting, err := ioutil.ReadFile("../stats/ting.mp3")
	//if err != nil {
	//	log.Fatalln(err)
	//}

	os.Getenv("")

	var resp models.TtsResp
	err := urllib.Post("http://tts_api.mechat.live/google_tts").SetJsonObject(models.TtsModel{
		Text: "警告",
		Lang: "ja",
	}).FromJsonByCode(&resp, 200)
	if err != nil {
		log.Fatalln(err)
	}

	code, bt, err := urllib.Get(resp.Url).RandDisguisedIP().RandUserAgent().ByteOriginal()
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
	create.Write(bt)

}
