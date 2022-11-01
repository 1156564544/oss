package objects

import (
	"crypto/sha256"
	"io"
	"encoding/base64"
)

func calculateHash(r io.Reader)string {
	h := sha256.New()
	io.Copy(h, r)
	hash:= base64.StdEncoding.EncodeToString(h.Sum(nil))
	return hash[:len(hash)-1]
}