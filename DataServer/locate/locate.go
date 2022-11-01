package locate

import (
	"os"
	"redisTool"
	"sync"
	"path/filepath"
	"fmt"
)

var mu sync.Mutex
var cache = make(map[string]bool)

func Add(hash string){
	mu.Lock()
	cache[hash] = true
	mu.Unlock()
}

func Delete(hash string){
	mu.Lock()
	delete(cache, hash)
	mu.Unlock()
}

func Locate(hash string) bool {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println(cache,hash)
	if _, ok := cache[hash]; ok {
		fmt.Println("true")
		return true
	}
	return false
}

func StartLocate(ip string) {
	names := redisTool.SubMessage("dataServers")
	for name := range names {
		if Locate(name) {
			redisTool.PushMessage(name, ip)
		}
	}
}

func init() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		hash := filepath.Base(files[i])
		cache[hash] = true
	}
}
