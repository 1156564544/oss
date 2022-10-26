package heartbeat

import (
	"redisTool"
	"time"
)

func StartHeartbeat(ip string) {
	for true {
		// 每隔5秒钟向apiServers广播自己的IP地址
		redisTool.PubMessage("apiServers", ip)
		time.Sleep(5 * time.Second)
	}
}
