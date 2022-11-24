package manage

import (
	"Proxy/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"redisTool"
	"strings"
)

// 新增用户
func AddUser(w http.ResponseWriter, r *http.Request) {
	token:=r.Header.Get("Authorization")
	if token==""{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if ok,err:=redisTool.SetExist(token);err!=nil||!ok{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	msg,err:=utils.DecryptByAes(token)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user utils.Users
	err=json.Unmarshal(msg,&user)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user.Isroot!=1{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	b,err:=ioutil.ReadAll(r.Body)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var newUser utils.Users
	err=json.Unmarshal(b,&newUser)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err=utils.InsertIntoDB(&newUser)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}


// 删除用户
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method!=http.MethodDelete{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	token:=r.Header.Get("Authorization")
	if token==""{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if ok,err:=redisTool.SetExist(token);err!=nil||!ok{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	username:=strings.Split(r.URL.EscapedPath(), "/")[2]
	err:=utils.DeleteFromDB(username)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// 修改用户
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method!=http.MethodPut{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	token:=r.Header.Get("Authorization")
	if token==""{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if ok,err:=redisTool.SetExist(token);err!=nil||!ok{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	username:=strings.Split(r.URL.EscapedPath(), "/")[2]
	b,err:=ioutil.ReadAll(r.Body)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user utils.Users
	err=json.Unmarshal(b,&user)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err=utils.UpdateDB(username,&user)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}