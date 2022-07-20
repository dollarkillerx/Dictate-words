package server

import (
	"github.com/dollarkillerx/Dictate-words/pkg/models"
	"github.com/dollarkillerx/async_utils"
	"github.com/dollarkillerx/processes"
	"github.com/dollarkillerx/urllib"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"time"
)

type GenerateTTSRequest struct {
	Lang        string `json:"lang" binding:"required"`
	Text        string `json:"text" binding:"required"`
	RepeatTimes int    `json:"repeat_times" binding:"required"`
}

func (s *Server) generateTTS(ctx *gin.Context) {
	var payload GenerateTTSRequest
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}
	if payload.RepeatTimes > 3 {
		ctx.JSON(400, "重複次數 最大3次")
		return
	}

	var words []Word
	rps := strings.Split(payload.Text, "\n")
	for _, v := range rps {
		vv := strings.TrimSpace(v)
		if vv != "" {
			words = append(words, Word{Word: vv, ID: xid.New().String()})
		}
	}

	if len(words) > 100 {
		ctx.JSON(400, "每次生成最大100行")
		return
	}

	cXid := xid.New().String()

	var over = make(chan struct{})

	poolFunc := async_utils.NewSinglePool(1, func() {
		close(over)
	})

	for i := range words {
		idx := i
		xp := xid.New().String()
		poolFunc.Send(func() error {
			px, err := sendPX(words[idx].Word, payload.Lang, cXid, xp)
			if err != nil {
				log.Println(err)
				return err
			}

			words[idx].FileName = px
			return nil
		})
	}

	poolFunc.Over()
	<-over

	err = poolFunc.Error()
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	// ffmpeg -i "concat:./stats/start.mp3|./test/xxx.mp3|./stats/ting.mp3" -acodec copy output.mp3 -y
	concat := "stats/start.mp3|"

	for _, v := range words {
		for i := 0; i < payload.RepeatTimes; i++ {
			concat += fmt.Sprintf("%s|", v.FileName)
			concat += "stats/ting.mp3|"
		}
	}

	concat += "stats/ting.mp3"

	cm := fmt.Sprintf(`ffmpeg -i "%s" -acodec copy %s.mp3 -y`, concat, cXid)
	log.Println(cm)
	_, err = processes.RunCommand(cm)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, err.Error())
		return
	}
	ctx.JSON(200, gin.H{
		"id": cXid,
	})
}

func sendPX(text string, lang string, prefix string, xp string) (string, error) {
	var resp models.TtsResp
	err := urllib.Post("http://tts_api.mechat.live/google_tts").SetJsonObject(models.TtsModel{
		Text: text,
		Lang: lang,
	}).FromJsonByCode(&resp, 200)
	if err != nil {
		return "", err
	}

	code, bt, err := urllib.Get(resp.Url).SetTimeout(time.Second * 10).RandDisguisedIP().RandUserAgent().ByteOriginal()
	if err != nil {
		return "", err
	}

	if code != 200 {
		return "", fmt.Errorf("%s", bt)
	}

	filename := path.Join("stats", "temporary", fmt.Sprintf("%s_%s.mp3", prefix, xp))
	err = ioutil.WriteFile(filename, bt, 00777)
	if err != nil {
		return "", err
	}

	return filename, nil
}
