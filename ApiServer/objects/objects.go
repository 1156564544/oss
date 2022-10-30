package objects

import (
	"ApiServer/heartbeat"
	"ApiServer/locate"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"es"
	"httpTool"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodGet && r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	switch r.Method {
	case http.MethodGet:
		object := strings.Split(r.URL.EscapedPath(), "/")[2]
		versionId := r.URL.Query()["version"]
		version := 0
		var err error
		if len(versionId) != 0 {
			// URL 指定了version
			version, err = strconv.Atoi(versionId[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		meta, err := es.GetMetadata(object, version)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if meta.Hash == "" {
			// 该对象已被删除
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// 从dataServers获取对象数据
		ip := locate.Locate(meta.Hash)
		if len(ip) == 0 {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("%v is not exist!\n", object)
			return
		}
		resp, err := http.Get("http://localhost" + ip + "/objects/" + meta.Hash)
		if err != nil {
			log.Println(err.Error())
		}
		defer resp.Body.Close()
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	case http.MethodPut:
		// 从URL中获取对象的hash和size
		hash := httpTool.GetHashFromHeader(r.Header)
		if hash == "" {
			log.Println("Hash is missing!")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		size := httpTool.GetSizeFromHeader(r.Header)

		// 写入对象数据到dataServers
		if locate.Locate(hash) == "" {
			ip := heartbeat.RandomChooseDataServers(1)[0]
			log.Println(ip)
			url := "http://localhost" + ip + "/objects/" + hash
			req, err := http.NewRequest("PUT", url, r.Body)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			res, _ := http.DefaultClient.Do(req)
			defer res.Body.Close()
		}else{
			log.Printf("%v is exist!\n", hash)
		}

		// 写入元数据到es
		object := strings.Split(r.URL.EscapedPath(), "/")[2]
		err := es.AddVersion(object, hash, size)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case http.MethodDelete:
		object := strings.Split(r.URL.EscapedPath(), "/")[2]
		meta, err := es.SearchLatestVersion(object)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = es.PutMetadata(meta.Name, meta.Version+1, 0, "")
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
