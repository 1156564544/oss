package versions

import (
	"encoding/json"
	"net/http"
	"es"
	"log"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// Get the object name from the request.
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	// Get the object version from the request.
	from, size := 0, 1000
	for true {
		metas, err := es.SearchAllVersions(name, from, size)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, meta := range metas {
			b, _ := json.Marshal(meta)
			w.Write(b)
			w.Write([]byte("\n"))
		}
		if len(metas) < size {
			break
		}
		from += size
	}
}
