package httpTool

import (
	"net/http"
	"strconv"
)

func GetHashFromHeader(h http.Header)(digest string){
	digest =h.Get("Digest")
	if len(digest)<=8||digest[:8]!="SHA-256="{
		digest=""
	}else{
		digest=digest[8:]
	}
	return
}

func GetSizeFromHeader(h http.Header)int64{
	content_length:=h.Get("Content-Length")
	size,_:=strconv.ParseInt(content_length,10,64)
	return size
}