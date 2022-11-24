package objects

import (
	"ApiServer/locate"
	"crypto/sha256"
	"io"
	"encoding/base64"
	"net/url"
	"net/http"
	"log"
	"encoding/json"

	"rs"
	"db"
	"aes"
	"redisTool"
)

// 计算hash值
func calculateHash(r io.Reader)string {
	h := sha256.New()
	io.Copy(h, r)
	hash:= base64.StdEncoding.EncodeToString(h.Sum(nil))
	return url.PathEscape(hash)
}

// 判断object是否存在
func Exist(hash string)bool{
	return len(locate.Locate(hash)) >= rs.NUM_DATA_SHARES
}

// 计算单个chunk的大小
func getPersharedSize(size int64)int64{
	return (size + int64(rs.NUM_DATA_SHARES-1)) / int64(rs.NUM_DATA_SHARES)
}

// 计算一个round读写的数据量
func getRoundSize()int{
	return rs.CHUNK_SIZE*(rs.NUM_DATA_SHARES+rs.NUM_PARITY_SHARES)
}

// 验证是否具有读权限
func checkReadPermission(r *http.Request)bool{
	authorization:=r.Header.Get("Authorization")
	if authorization==""{
		log.Println("Authorization is missing!")
		return false
	}
	if ok,err:=redisTool.SetExist(authorization);err!=nil||!ok{
		log.Println("Authorization is error!")
		return false
	}
	var user db.Users
	b,err:=aes.DecryptByAes(authorization)
	if err!=nil{
		log.Println("DecryptByAes error:",err.Error())
		return false
	}
	err=json.Unmarshal(b,&user)
	if err!=nil{
		log.Println("json.Unmarshal error:",err.Error())
		return false
	}
	if user.Isread!=1{
		log.Println("user.Isread!=1")
		return false
	}
	return true
}

// 验证是否具有写权限
func checkWritePermission(r *http.Request)bool{
	authorization:=r.Header.Get("Authorization")
	if authorization==""{
		log.Println("Authorization is missing!")
		return false
	}
	if ok,err:=redisTool.SetExist(authorization);err!=nil||!ok{
		log.Println("Authorization is error!")
		return false
	}
	var user db.Users
	b,err:=aes.DecryptByAes(authorization)
	if err!=nil{
		log.Println("DecryptByAes error:",err.Error())
		return false
	}
	err=json.Unmarshal(b,&user)
	if err!=nil{
		log.Println("json.Unmarshal error:",err.Error())
		return false
	}
	if user.Iswrite!=1{
		log.Println("user.Iswrite!=1")
		return false
	}
	return true
}
