package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// GenMD5 GenMD5
func GenMD5(src interface{}) string {
	var bts []byte
	switch r := src.(type) {
	case string:
		bts = []byte(r)
	case []byte:
		bts = r
	}

	sum := md5.New().Sum(bts)

	return hex.EncodeToString(sum)
}
