package Client

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"strings"
	"net/url"
	"log"
)


func CalculateHash(s string) string {
	r := strings.NewReader(s)
	h := sha256.New()
	io.Copy(h, r)
	hash := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return url.PathEscape(hash)
}

func CalculateHashWithReader(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)
	hash := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return url.PathEscape(hash)
}

func CheckErr(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
