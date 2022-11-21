package redisTool

import (
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	key:= "test"
	SetAdd(key)

	time.Sleep((setDuration-1)*time.Second)
	ok,err:=SetExist(key)
	if err!=nil{
		t.Error(err)
	}
	// 此时还没有过期，应该找得到
	if !ok {
		t.Error("SetExist error")
	}
	
	time.Sleep(2*time.Second)
	ok,err=SetExist(key)
	if err!=nil{
		t.Error(err)
	}
	// 此时过期了，应该找不到
	if ok {
		t.Error("SetExist error with expire")
	}
}