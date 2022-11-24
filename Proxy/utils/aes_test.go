package utils

import (
	"testing"
)

func Test_aes(t *testing.T) {
	msg:=[]byte("hello world")
	enc,err:=AesEncrypt(msg,PwdKey)
	if err!=nil{
		t.Error(err)
	}
	dec,err:=AesDecrypt(enc,PwdKey)
	if err!=nil{
		t.Error(err)
	}
	if string(dec)!=string(msg){
		t.Error("decrypt fail")
	}
}