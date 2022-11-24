package main

import (
	"fmt"
	"io"
	"crypto/tls"
	"log"
	"net/http"
	"strconv"
	"strings"
	"bytes"
	"os"
	utils "Client/utils"
	"encoding/json"

	"redisTool"
)

type Users struct{
	Name string
	Password string
	Isroot int
	Isread int
	Iswrite int
}

// 通过Post方法获取token
func PostHeader(url string, headers map[string]string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("post failed with status code:", resp.StatusCode)
		return "", fmt.Errorf("post failed with status code:%d", resp.StatusCode)
	}
	token := resp.Header.Get("location")
	return token, nil
}

// 删除对象
func delete(object string, authorization string) error {
	url := "http://localhost:10000/objects/" + object
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, strings.NewReader(""))
	req.Header.Set("Authorization", authorization)
	if utils.CheckErr(err) {
		return err
	}
	resp, err := client.Do(req)
	if utils.CheckErr(err) {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("delete failed with status code:", resp.StatusCode)
		return fmt.Errorf("delete failed with status code:%d", resp.StatusCode)
	}
	return nil
}


// post方法用来通过段点上传的方式上传大对象
// 1.调用PostHeader方法获取token
// 2.调用putFromToken方法进行段点上传
func post(object string, filepath string, authorization string) error {
	stat,err:=os.Stat(filepath)
	if os.IsNotExist(err){
		log.Println("file not exist")
		return fmt.Errorf("file not exist")
	}
	size:=stat.Size()
	f,err:=os.Open(filepath)
	if utils.CheckErr(err) {
		return err
	}
	defer f.Close()
	hash := utils.CalculateHashWithReader(f)
	url := "http://localhost:10000/objects/" + object
	headers := make(map[string]string)
	headers["Digest"] = "SHA-256=" + hash
	headers["length"] = strconv.FormatInt(size, 10)
	headers["Authorization"] = authorization
	token, err := PostHeader(url, headers)
	if utils.CheckErr(err) {
		return err
	}
	// tf,_:=os.Create(object+".token")
	// _,err=tf.Write([]byte(token))
	// if utils.CheckErr(err) {
	// 	return err
	// }
	// defer tf.Close()
	redisTool.AddKeyValue(object, token)
	err=putFromToken(token, filepath)
	if utils.CheckErr(err) {
		return err
	}
	// os.Remove(filepath+".token")
	redisTool.DelKeyValue(object)
	return nil
}


// 通过head方法查询当前上传的进度
func head(token string)int64{
	client := &http.Client{}
	url := "http://localhost:10000" + token
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if utils.CheckErr(err) {
		return 0
	}
	resp, err := client.Do(req)
	if utils.CheckErr(err) {
		return 0
	}
	defer resp.Body.Close()
	current:=resp.Header.Get("content-length")
	size,err:=strconv.ParseInt(current, 10, 64)
	if utils.CheckErr(err) {
		return 0
	}
	return size
}

// 通过token进行段点上传
// 返回nil意味着上传完成
func putFromToken(token string, filepath string) error {
	offset:=head(token)
	stat,err:=os.Stat(filepath)
	if os.IsNotExist(err){
		log.Println("file not exist")
		return fmt.Errorf("file not exist")
	}
	size:=stat.Size()
	if offset==size{
		return nil
	}
	f,err:=os.Open(filepath)
	if utils.CheckErr(err) {
		return err
	}
	defer f.Close()
	_,err=f.Seek(offset,io.SeekStart)
	if utils.CheckErr(err) {
		return err
	}

	client := &http.Client{}
	url := "http://localhost:10000" + token
	req, err := http.NewRequest("PUT", url, f)
	if utils.CheckErr(err){
		return err
	}
	headers := make(map[string]string)
	headers["Range"] = "bytes=" + strconv.FormatInt(offset, 10) + "_" +  " /"
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	resp, err := client.Do(req)
	if utils.CheckErr(err){
		return err
	}
	if resp.StatusCode!=http.StatusOK{
		log.Println("put failed with status code:",resp.StatusCode)
		return fmt.Errorf("put failed with status code:%d",resp.StatusCode)
	}
	return nil
}


// 上传对象
// object:对象名，filepath:上传的文件路径
func put(object string, filepath string,authorization string) error {
	stat,err:=os.Stat(filepath)
	if os.IsNotExist(err){
		log.Println("file not exist")
		return fmt.Errorf("file not exist")
	}

	size:=stat.Size()
	if size>10{
		// 对于100MVB以上的文件，需要采取段点上传的方式
		return post(object,filepath,authorization)
	}

	f,err:=os.Open(filepath)
	if utils.CheckErr(err) {
		return err
	}
	defer f.Close()

	msg:=make([]byte,size)
	_,err=f.Read(msg)
	if utils.CheckErr(err) {
		return err
	}
	hash:=utils.CalculateHash(string(msg))

	client := &http.Client{}
	url := "http://localhost:10000/objects/" + object
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(msg)))
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["length"] = strconv.Itoa(len(msg))
	headers["Digest"] = "SHA-256=" + hash
	headers["Authorization"] = authorization
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	resp, err := client.Do(req)
	if utils.CheckErr(err) {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode!=http.StatusOK{
		log.Println("put failed with status code:",resp.StatusCode)
		return fmt.Errorf("put failed with status code:%d",resp.StatusCode)
	}
	return nil
}


// 下载对象
// object:对象名，filepath:下载到的文件路径，offset:下载的偏移量
func get(object string,filepath string,offset int64,authorization string) error {
	url := "http://localhost:10000/objects/" + object
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	req.Header.Set("Authorization",authorization)
	if utils.CheckErr(err) {
		return err
	}
	if offset > 0 {
		req.Header.Set("Range", "bytes=" + strconv.FormatInt(offset, 10) + "_" )
	}
	resp, err := client.Do(req)
	if utils.CheckErr(err) {
		return err
	}
	if resp.StatusCode!=http.StatusOK{
		log.Println("get failed with status code:",resp.StatusCode)
		return fmt.Errorf("get failed with status code:%d",resp.StatusCode)
	}
	if offset > 0 {
		resp.Header.Set("Range", "bytes=" + strconv.FormatInt(offset, 10) + "_" + strconv.FormatInt(offset+int64(resp.ContentLength)-1, 10) + "/" + strconv.FormatInt(int64(resp.ContentLength), 10))
	}
	defer resp.Body.Close()
	_,err=os.Stat(filepath)
	if os.IsNotExist(err) && offset!=0{
		log.Println("file not exist but offset is not 0")
		return fmt.Errorf("file not exist but offset is not 0")
	}
	f,err:=os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, os.ModeAppend|os.ModePerm)
	f.Seek(offset, io.SeekStart)
	if utils.CheckErr(err) {
		return err
	}
	_,err=io.Copy(f, resp.Body)
	if utils.CheckErr(err) {
		return err
	}
	return nil
}

func login(username string, password string) (string,error){
	// 不验证服务端证书
    tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
    }

    client := &http.Client{Transport: tr}

    req, err := http.NewRequest(http.MethodGet, "https://localhost:9999/login", http.NoBody)
	if utils.CheckErr(err){
		return "",err
	}
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if utils.CheckErr(err){
		return "",err
	}
	if resp.StatusCode!=http.StatusOK{
		log.Println("login failed with status code:",resp.StatusCode)
		return "",fmt.Errorf("login failed with status code:%d",resp.StatusCode)
	}
	token:=resp.Header.Get("Authorization")
	return token,nil
}

func addUser(newUser *Users,token string)bool{
	// 不验证服务端证书
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	b,_:=json.Marshal(&newUser)
	req, err := http.NewRequest(http.MethodPost, "https://localhost:9999/addUser", bytes.NewReader(b))
	if utils.CheckErr(err){
		return false
	}
	req.Header.Add("Authorization", token)
	resp, err := client.Do(req)
	if utils.CheckErr(err){
		return false
	}
	if resp.StatusCode!=http.StatusOK{
		log.Println("add user failed with status code:",resp.StatusCode)
		return false
	}
	return true
}

func deleteUser(username string,token string)bool{
	// 不验证服务端证书
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(http.MethodDelete, "https://localhost:9999/deleteUser/"+username, http.NoBody)
	if utils.CheckErr(err){
		return false
	}
	req.Header.Add("Authorization", token)
	resp, err := client.Do(req)
	if utils.CheckErr(err){
		return false
	}
	if resp.StatusCode!=http.StatusOK{
		log.Println("delete user failed with status code:",resp.StatusCode)
		return false
	}
	return true
}

func updateUser(username string,token string,newUser *Users)bool{
	// 不验证服务端证书
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	b,_:=json.Marshal(&newUser)
	req, err := http.NewRequest(http.MethodPut, "https://localhost:9999/updateUser/"+username, bytes.NewReader(b))
	if utils.CheckErr(err){
		return false
	}
	req.Header.Add("Authorization", token)
	resp, err := client.Do(req)
	if utils.CheckErr(err){
		return false
	}
	if resp.StatusCode!=http.StatusOK{
		log.Println("update user failed with status code:",resp.StatusCode)
		return false
	}
	return true
}

func main() {
	// object := "test7_29"
	// filepath:="./"+object
	// var offset int64 = 4


	// msg := "this is the object test7_25"
	// token, err := post(object, msg)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }
	// fmt.Println(token)
	// err = putFromToken(token, 0, 8,msg)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }
	// err = putFromToken(token, 8, int64(len(msg))-8,msg)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }
	// put(object,filepath)
	// get(object,"./"+object,offset)

	// authorization,_:=login("root", "123456")
	 authorization:="zIWveE2D9wkoDcn5yCCnRHAJ8nBSnbfxzdldAtSkk5z71zRAB37leNRST+G0HhJIC2UO6EpKCcwxtHydt6F7mp/HGDY6KCMohEjxsNT2ZcM="

	// if err!=nil{
	// 	log.Println(err.Error())
	// 	return
	// }
	// fmt.Println(authorization)
	// put("test8_1","./test8_1",authorization)
	get("test8_1","./test8_1",0,authorization)
}
