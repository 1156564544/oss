package objects

import (
	"DataServer/utils"
	"net/http"
	"os"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	// 基于gzip解压缩
	err := utils.UnGzip(string(os.Getenv("STORAGE_ROOT")+"/objects/"+strings.Split(r.URL.EscapedPath(), "/")[2]), &w)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 没有进行压缩时直接读取文件
	// f, err := os.Open(os.Getenv("STORAGE_ROOT") + "/objects/" + strings.Split(r.URL.EscapedPath(), "/")[2])
	// if err != nil {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	// defer f.Close()
	// io.Copy(w, f)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodGet {
		get(w, r)
	}
}
