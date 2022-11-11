package locate

import (
	"encoding/json"
	"log"
	"net/http"
	"redisTool"
	"strconv"
	"strings"
	"time"

	"rs"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	ips := Locate(object)
	log.Println(ips)
	if len(ips) < rs.NUM_DATA_SHARES {
		w.WriteHeader(http.StatusNotFound)
	} else {
		b, _ := json.Marshal(ips)
		w.Write(b)
	}
}

func Locate(object string) (ips map[int]string) {
	ips = make(map[int]string)
	redisTool.PubMessage("dataServers", object)
	now := time.Now()
	for now.Add(1 * time.Second).After(time.Now()) {
		// ip:<ip>_<id of Shard>
		msg := redisTool.PopMessage(object)
		if msg != "" {
			ip:=strings.Split(msg, "_")[0]
			id,_:=strconv.Atoi(strings.Split(msg, "_")[1])
			ips[id]=ip
		}
		// 已经接收到了所有的分片             
		if len(ips)==rs.NUM_DATA_SHARES+rs.NUM_PARITY_SHARES{
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
	return 
}
