package objects

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type putStream struct {
	uuid       string
	dataServer string
}

func CreatePutStream(server string, hash string, size int64) (*putStream, error) {
	// 向dataServer的temp接口发送post请求，同时捎带hash和size
	resp, err := http.PostForm("http://localhost"+server+"/temp/"+hash,
		url.Values{"Size": {fmt.Sprintf("%d", size)}})
	if err != nil {
		return nil, fmt.Errorf("Post to dataServer error: %v", err.Error())
	}
	defer resp.Body.Close()
	// 从dataServer的响应中获取uuid
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,fmt.Errorf("Read uuid from dataServer error: %v", err.Error())
	}
	uuid := string(body)
	// log.Println("uuid:", uuid)
	// log.Println("server:", server)
	// 根据dataServer的IP地址和uuid构造putStream对象
	return &putStream{uuid,server}, nil
}

func (stream *putStream) Write(p []byte) (n int, err error) {
	req,err:=http.NewRequest(http.MethodPatch, "http://localhost"+stream.dataServer+"/temp/"+stream.uuid, strings.NewReader(string(p)))
	if err!=nil{
		return 0,err
	}
	client:=http.Client{}
	resp,err:=client.Do(req)
	if err!=nil{
		return 0,err
	}
	defer resp.Body.Close()
	if resp.StatusCode!=http.StatusOK{
		return 0,fmt.Errorf("dataServer write error with status code %d",resp.StatusCode)
	}
	return len(p),nil
}

func (stream *putStream)commit(iscommit bool) error{
	method:=http.MethodPut
	if !iscommit{
		method=http.MethodDelete
	}
	req,_:=http.NewRequest(method, "http://localhost"+stream.dataServer+"/temp/"+stream.uuid,nil)
	client := http.Client{}
	resp,err:=client.Do(req)
	if err!=nil{
		return err
	}
	if resp.StatusCode!=http.StatusOK{
		return fmt.Errorf("dataServer commit error with status code %d",resp.StatusCode)
	}
	return nil
}
