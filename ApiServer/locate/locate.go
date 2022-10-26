package locate

import (
	"encoding/json"
	"log"
	"net/http"
	"redisTool"
	"strings"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	ip := Locate(object)
	log.Println(ip)
	if len(ip) == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		b, _ := json.Marshal(ip)
		w.Write(b)
	}
}

func Locate(object string) string {
	redisTool.PubMessage("dataServers", object)
	now := time.Now()
	for now.Add(1 * time.Second).After(time.Now()) {
		ip := redisTool.PopMessage(object)
		if len(ip) != 0 {
			return ip
		}
		time.Sleep(1 * time.Millisecond)
	}
	return ""
}
