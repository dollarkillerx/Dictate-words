package server

import (
	"github.com/dollarkillerx/Dictate-words/pkg/models"
	"github.com/dollarkillerx/Dictate-words/pkg/utils"
	"github.com/dollarkillerx/async_utils"
	"github.com/dollarkillerx/processes"
	"github.com/dollarkillerx/urllib"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"math/rand"

	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type GenerateTTSRequest struct {
	Lang        string `json:"lang" binding:"required"`
	Text        string `json:"text" binding:"required"`
	RepeatTimes int    `json:"repeat_times" binding:"required"`
	PlayOrder   string `json:"play_order"`
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

	get, ok := s.cacheGet(payload.Lang+payload.PlayOrder, payload.Text, payload.RepeatTimes)
	if ok {
		ctx.JSON(200, gin.H{
			"id": get,
		})
		return
	}

	//del
	defer func() {
		for _, v := range words {
			if v.FileName != "" {
				os.Remove(v.FileName)
			}
		}
	}()

	cXid := xid.New().String()

	var over = make(chan struct{})

	poolFunc := async_utils.NewSinglePool(3, func() {
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

	// https://blog.csdn.net/tian2342/article/details/99303883

	err = poolFunc.Error()
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	// ffmpeg -i "concat:./stats/start.mp3|./test/xxx.mp3|./stats/ting.mp3" -acodec copy output.mp3 -y
	concat := "concat:stats/start.mp3|"
	var newWord string

	if payload.PlayOrder == "random" {
		var newWords []Word
		to := len(words)
		for i := to - 1; i > 0; i-- {
			rand.Seed(time.Now().UnixNano())
			ri := rand.Intn(i)
			newWords = append(newWords, words[ri])
			newWord += fmt.Sprintf("%s\n", words[ri].Word)
			if ri == len(words)-1 {
				words = append(words[:ri])
			} else {
				words = append(words[:ri], words[ri+1:]...)
			}
		}
		words = newWords
	}

	for _, v := range words {
		for i := 0; i < payload.RepeatTimes; i++ {
			concat += fmt.Sprintf("%s|", v.FileName)
			concat += "stats/ting.mp3|"
		}
	}

	concat += "stats/ting.mp3"

	cm := fmt.Sprintf(`ffmpeg -i "%s" -acodec copy stats/temporary/%s.mp3 -y`, concat, cXid)
	log.Println(cm)
	_, err = processes.RunCommand(cm)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, err.Error())
		return
	}

	s.cacheSet(payload.Lang+payload.PlayOrder, payload.Text, payload.RepeatTimes, cXid, fmt.Sprintf("stats/temporary/%s.mp3", cXid))

	ctx.JSON(200, gin.H{
		"id":   cXid,
		"word": newWord,
	})
}

func (s *Server) cacheGet(lang string, words string, repeatTimes int) (id string, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := utils.GenMD5(fmt.Sprintf("%s_%s_%d", lang, words, repeatTimes))
	cache, ok := s.cache[key]
	if ok {
		if cache.Expiration > time.Now().Unix() {
			return cache.ID, true
		}
	}

	return "", false
}

func (s *Server) cacheSet(lang string, words string, repeatTimes int, id string, filepath string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := utils.GenMD5(fmt.Sprintf("%s_%s_%d", lang, words, repeatTimes))
	cache, ok := s.cache[key]
	if ok {
		os.Remove(cache.Filepath)
	}

	s.cache[key] = WordCache{
		ID:         id,
		Expiration: time.Now().Add(time.Hour).Unix(),
		Filepath:   filepath,
	}
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

func (s *Server) getDownloadTTSPath(ttsID string) (path string, ex bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range s.cache {
		if v.ID == ttsID {
			if v.Expiration > time.Now().Unix() {
				return v.Filepath, true
			}
		}
	}

	return "", false
}

func (s *Server) downloadTTS(ctx *gin.Context) {
	ttsID := ctx.Param("tts_id")
	if ttsID == "" {
		ctx.String(400, "HTTP 400 資源錯誤")
		return
	}

	ttsPath, ex := s.getDownloadTTSPath(ttsID)
	if !ex {
		ctx.String(400, "HTTP 4001 資源過期")
		return
	}

	open, err := os.Open(ttsPath)
	if err != nil {
		log.Println(err)
		ctx.String(400, "HTTP 4001 資源過期")
		return
	}
	ctx.Header("Content-Type", "audio/mpeg")
	ctx.Header("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s.mp3", ttsID)) // 用来指定下载下来的文件名
	io.Copy(ctx.Writer, open)
}
