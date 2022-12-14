package heartbeat

import (
	"math/rand"
	"redisTool"
	"sync"
	"time"
)

var dataServers map[string]time.Time
var mu sync.Mutex

func init() {
	dataServers = make(map[string]time.Time)
}

func ListenHeartbeat() {
	go removeDataServers()
	ips := redisTool.SubMessage("apiServers")
	for ip := range ips {
		mu.Lock()
		dataServers[ip] = time.Now()
		mu.Unlock()
	}
}

func removeDataServers() {
	for true {
		time.Sleep(10 * time.Second)
		mu.Lock()
		for servers, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, servers)
			}
		}
		mu.Unlock()
	}
}

func GetDataServers() []string {
	res := []string{}
	for servers, _ := range dataServers {
		res = append(res, servers)
	}
	return res
}

func RandomChooseDataServers(n int) []string {
	mu.Lock()
	defer mu.Unlock()
	dataServerIps := GetDataServers()
	if len(dataServerIps) < n {
		n = len(dataServerIps)
	}
	idx := 0
	servers := make([]string, n)
	indexs := make(map[int]bool)
	for idx < n {
		index := rand.Intn(n)
		if !indexs[index] {
			indexs[index] = true
			servers[idx] = dataServerIps[index]
			idx += 1
		}
	}
	return servers
}

func randomChooseDataServersWithExclude(dataServerIps []string,n int,exclude map[int]string)[]string{
	mu.Lock()
	defer mu.Unlock()
	exist:=make(map[string]bool)
	for _,ip:=range exclude{
		exist[ip]=true
	}
	servers:=make([]string,0)
	for len(servers)+len(exclude)<n{
		index:=rand.Intn(len(dataServerIps))
		if !exist[dataServerIps[index]]{
			servers=append(servers,dataServerIps[index])
			exist[dataServerIps[index]]=true
		}
	}
	return servers
}

func RandomChooseDataServersWithExclude(n int, exclude map[int]string) []string {
	return randomChooseDataServersWithExclude(GetDataServers(),n,exclude)
}