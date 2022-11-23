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

func TestKeyValue(t *testing.T) {
	key:= "key"
	value:= "value"
	err:=AddKeyValue(key,value)
	if err!=nil{
		t.Error(err)
	}

	v,err:=GetKeyValue(key)
	if err!=nil{
		t.Error(err)
	}
	if v!=value{
		t.Error("GetKeyValue error")
	}
	DelKeyValue(key)
	v,err=GetKeyValue(key)

	if v!=""{
		t.Error("DelKeyValue error")
	}
}