package objects

import (
	"ApiServer/heartbeat"
	"ApiServer/locate"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodGet {
		object := strings.Split(r.URL.EscapedPath(), "/")[2]
		ip := locate.Locate(object)
		if len(ip) == 0 {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("%v is not exist!\n", object)
			return
		}
		resp, err := http.Get("http://localhost" + ip + "/objects/" + object)
		if err != nil {
			log.Println("err")
		}
		defer resp.Body.Close()
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		object := strings.Split(r.URL.EscapedPath(), "/")[2]
		ip := heartbeat.RandomChooseDataServers(1)[0]
		fmt.Println(ip)
		url := "http://localhost" + ip + "/objects/" + object
		req, err := http.NewRequest("PUT", url, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()
	}
}
