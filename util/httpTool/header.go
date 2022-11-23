package httpTool

import (
	"net/http"
	"strconv"
	"strings"
	"log"
)

// 从http头中读取object的hash
func GetHashFromHeader(h http.Header) (digest string) {
	digest = h.Get("Digest")
	if len(digest) <= 8 || digest[:8] != "SHA-256=" {
		digest = ""
	} else {
		digest = digest[8:]
	}
	return
}

// 从http头中读取object的size
func GetSizeFromHeader(h http.Header) int64 {
	content_length := h.Get("length")
	size, _ := strconv.ParseInt(content_length, 10, 64)
	return size
}

// 从http头中读取object的offset
func GetOffsetFromHeader(h http.Header) int64 {
	range_content := h.Get("Range")
	// fmt.Println(range_content)
	if len(range_content) <= 6 || range_content[:6] != "bytes=" {
		log.Println("Range is missing: ", range_content)
		return 0
	}
	range_content = range_content[6:]
	offset, e := strconv.ParseInt(strings.Split(range_content, "_")[0], 10, 64)
	if e != nil {
		log.Println(e)
		return 0
	}
	return offset
}
