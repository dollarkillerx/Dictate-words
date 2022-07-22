package googletts

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTTSURL(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		str  string
		lang string
	}{
		{"a", "en"},
		{"b", "en"},
		{"foo", "en"},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.", "en"},
		//{"こんにちは、世界。", "ja"},
	}

	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		u, err := GetTTSURL(tt.str, tt.lang)
		a.NoError(err, target)
		a.NotEmpty(u, err, target)
		a.Contains(u, "https://translate.google.com/translate_tts?", target)
	}
}

func TestRand(t *testing.T) {
	c := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var newC []int
	for i := len(c); i > 0; i-- {
		rand.Seed(time.Now().UnixNano())
		intn := rand.Intn(i)
		newC = append(newC, c[intn])
		if intn == len(c)-1 {
			c = append(c[:intn])
		} else {
			c = append(c[:intn], c[intn+1:]...)
		}
	}

	log.Println(newC)
}
