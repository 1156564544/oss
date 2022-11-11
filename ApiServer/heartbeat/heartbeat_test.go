package heartbeat

import (
	"testing"
)

func TestRandomChooseDataServersWithExclude(t *testing.T) {
	dataServers:=[]string{":10001",":10002",":10003",":10004",":10005",":10006",":10007",":10008"}
	exclude:=map[int]string{0:":10001",1:":10002",3:":10006",4:":10007"}
	exist:=make(map[string]bool)
	exist[":10001"]=true
	exist[":10002"]=true
	exist[":10006"]=true
	exist[":10007"]=true
	allServers:=make(map[string]bool)
	allServers[":10001"]=true
	allServers[":10002"]=true
	allServers[":10003"]=true
	allServers[":10004"]=true
	allServers[":10005"]=true
	allServers[":10006"]=true
	allServers[":10007"]=true
	allServers[":10008"]=true
	res:=randomChooseDataServersWithExclude(dataServers,6,exclude)
	if len(res)!=2{
		t.Errorf("Expected 2 servers,but got %d",len(res))
	}
	for i,server:=range res{
		if !allServers[server]{
			t.Errorf("server %s not in dataServers",server)
		}
		if exist[server]{
			t.Errorf("Expected server %s not in exclude,but got %s",server,res[i])
		}
	}
}