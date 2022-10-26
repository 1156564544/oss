package locate

import (
	"fmt"
	"os"
	"redisTool"
)

func locate(name string) bool {
	_, err := os.Stat(name)
	return err == nil || !os.IsNotExist(err)
}

func StartLocate(ip string) {
	names := redisTool.SubMessage("dataServers")
	for name := range names {
		fmt.Println(os.Getenv("STORAGE_ROOT") + "/objects/" + name)
		if locate(os.Getenv("STORAGE_ROOT") + "/objects/" + name) {
			redisTool.PushMessage(name, ip)
		}
	}
}
