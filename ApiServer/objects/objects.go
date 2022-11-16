package objects

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"ApiServer/heartbeat"
	"es"
	"httpTool"
	"rs"
)

func get(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	// log.Println(object)
	versionId := r.URL.Query()["version"]
	version := 0
	var err error
	if len(versionId) != 0 {
		// URL 指定了version
		version, err = strconv.Atoi(versionId[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	meta, err := es.GetMetadata(object, version)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if meta.Hash == "" {
		// 该对象已被删除
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// log.Println(meta.Hash,meta.Size)
	// 从dataServers获取对象数据
	stream,err:=CreateRSGetStream(meta.Hash,meta.Size)
	if err!=nil{
		log.Println("Get chunk error: ",err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	offset:=httpTool.GetOffsetFromHeader(r.Header)
	fmt.Println(offset)
	if offset!=0{
		// 从指定位置开始读取
		_,err=stream.Seek(offset,io.SeekCurrent)
	}
	_,err=io.Copy(w, stream)

	if err != nil {
		log.Println("Read error: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stream.Close()
}

func put(w http.ResponseWriter, r *http.Request) {
	// 从URL中获取对象的hash和size
	hash := httpTool.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("Hash is missing!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	size := httpTool.GetSizeFromHeader(r.Header)

	// 写入对象数据到dataServers
	if !Exist(hash)  {
		// 根据hash和size从dataServers获得uuid并创建put流
		dataServers:=heartbeat.RandomChooseDataServers(rs.NUM_PARITY_SHARES+rs.NUM_DATA_SHARES)
		stream,err:=CreateRSPutStream(dataServers,hash,size)
		if err!=nil{
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// 根据客户端的数据流计算内容的哈希值，并将客户端的数据流传输给dataServers
		reader:=io.TeeReader(r.Body, stream)
		calculateHash:=calculateHash(reader)
		// 如果计算出来的哈希值等于URL中的哈希值，则说明数据传输成功，让dataServers保存该数据，否则让dataServers删除该数据
		iscommit:=calculateHash==hash
		if !iscommit{
			log.Println("Hash is error!")
			stream.commit(iscommit)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
		err=stream.commit(iscommit)
		if err!=nil{
			log.Println("Streaam commit failed: "+err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		log.Printf("%v is exist!\n", hash)
	}

	// 写入元数据到es
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	err := es.AddVersion(object, url.PathEscape(hash), size)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	meta, err := es.SearchLatestVersion(object)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = es.PutMetadata(meta.Name, meta.Version+1, 0, "")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodGet && r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	switch r.Method {
	case http.MethodGet:
		get(w, r)
	case http.MethodPut:
		put(w, r)
	case http.MethodDelete:
		delete(w, r)
	}
}
