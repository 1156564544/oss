package temp

import (
	"crypto/sha256"
	"io"
	"encoding/base64"
	"net/url"
)

func calculateHash(r io.Reader)string {
	h := sha256.New()
	io.Copy(h, r)
	hash:= base64.StdEncoding.EncodeToString(h.Sum(nil))
	return url.PathEscape(hash)
}
