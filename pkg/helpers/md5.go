package helpers

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5Hash(s string) string {
	resp := md5.Sum([]byte(s))
	return hex.EncodeToString(resp[:])
}
