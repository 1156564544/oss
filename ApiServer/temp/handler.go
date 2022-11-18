package temp

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"ApiServer/objects"
)

func head(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, e := objects.GetRSResumablePutStreamFromToken(token)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("content-length", fmt.Sprintf("%d", current))
}


func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		put(w,r)
	case http.MethodHead:
		head(w,r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed) 
	}
}