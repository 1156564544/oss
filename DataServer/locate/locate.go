package locate

import (
	"fmt"
	"os"
	"path/filepath"
	"redisTool"
	"strconv"
	"strings"
	"sync"
)

var mu sync.Mutex
var cache = make(map[string]int)

// Params:hash--文件hash值,id--切片在文件中的id号
func Add(hash string,id int){
	mu.Lock()
	cache[hash] = id
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
	_, ok := cache[hash]
	return ok
}

func StartLocate(ip string) {
	names := redisTool.SubMessage("dataServers")
	for name := range names {
		hash:=strings.Split(name, ".")[0]
		if Locate(hash) {
			mu.Lock()
			id:=cache[hash]
			mu.Unlock()
			// 以<ip>_<id of shard>的格式向apiServers发送心跳
			// redisTool.PubMessage(name, ip+"_"+strconv.Itoa(id))
			fmt.Println(ip+"_"+strconv.Itoa(id))
			redisTool.PushMessage(hash, ip+"_"+strconv.Itoa(id))
		}
	}
}

func init() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		// 文件名格式：<hash of file>.<id>
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 2 {
			panic(files[i])
		}
		hash := file[0]
		id, e := strconv.Atoi(file[1])
		if e != nil {
			panic(e)
		}
		cache[hash] = id
	}
}
