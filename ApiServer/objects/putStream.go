package objects

import (
	"ApiServer/heartbeat"
	"errors"
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

func CreatePutStream(hash string, size string) (*putStream, error) {
	dataServers := heartbeat.RandomChooseDataServers(1)
	if len(dataServers) == 0 {
		return nil, errors.New("no dataServer")
	}
	server := dataServers[0]
	resp, err := http.PostForm("http://localhost"+server+"/temp/"+url.PathEscape(hash),
		url.Values{"Size": {size}})
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Post to dataServer error: %v", err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,fmt.Errorf("Read uuid from dataServer error: %v", err.Error())
	}
	uuid := string(body)
	fmt.Println(uuid)

	return &putStream{uuid,server}, nil
}

func (stream *putStream) Write(p []byte) (n int, err error) {
	fmt.Println("http://localhost"+stream.dataServer+"/temp/"+stream.uuid)
	req,err:=http.NewRequest(http.MethodPatch, "http://localhost"+stream.dataServer+"/temp/"+stream.uuid, strings.NewReader(string(p)))
	if err!=nil{
		return 0,err
	}
	client:=http.Client{}
	resp,err:=client.Do(req)
	defer resp.Body.Close()
	if err!=nil{
		return 0,err
	}
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
