package temp

import (
	"DataServer/locate"
	"DataServer/utils"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type TempInfo struct {
	Uuid string
	Hash string
	ID   int
	Size int64
}

// 处理来自apiServer的post请求，其中包括待上传的对象名（<hash>.<id of shard>)和单个shard的大小，dataServer生成uuid并作为post请求的响应返回给apiServer。
// dataServer将uuid作为文件名保存在temp目录下，文件内容为json格式的tempInfo，其中包括uuid、hash、id和size。
func post(w http.ResponseWriter, r *http.Request) {
	object_name := strings.Split(r.URL.EscapedPath(), "/")[2]
	// object_name: <hash>.<id of shard>
	hash := strings.Split(object_name, ".")[0]
	id, _ := strconv.Atoi(strings.Split(object_name, ".")[1])
	size, err := strconv.ParseInt(r.PostFormValue("Size"), 10, 64)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	output, _ := exec.Command("uuidgen").Output()
	uuid := strings.TrimSuffix(string(output), "\n")
	tempInfo := TempInfo{
		Uuid: uuid,
		Hash: hash,
		ID:   id,
		Size: size,
	}
	b, err := json.Marshal(tempInfo)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	f, err := os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid)
	defer f.Close()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = f.Write([]byte(b))
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(uuid))
	_, err = os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid + ".dat")
	if err != nil {
		log.Println(err.Error())
		os.Remove(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getTempInfo(uuid string) (TempInfo, error) {
	var tempInfo TempInfo
	f, err := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid)
	defer f.Close()
	if err != nil {
		return TempInfo{}, err
	}
	b, _ := ioutil.ReadAll(f)
	json.Unmarshal(b, &tempInfo)
	// decoder:=json.NewDecoder(f)
	// err=decoder.Decode(&tempInfo)
	if err != nil {
		return TempInfo{}, err
	}
	return tempInfo, nil
}

func patch(w http.ResponseWriter, r *http.Request) {
	// 从URL中获得uuid
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 根据uuid从磁盘读取tempInfo
	tempInfo, err := getTempInfo(uuid)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 创建用于保存临时数据的文件并将输入流中的数据写入临时文件
	// f, err := os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + tempInfo.Uuid + ".dat")
	f, err := os.OpenFile(os.Getenv("STORAGE_ROOT")+"/temp/"+tempInfo.Uuid+".dat", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		log.Println("Create file failed: " + err.Error())
		os.Remove(os.Getenv("STORAGE_ROOT") + "/temp/" + tempInfo.Uuid)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, r.Body)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func put(w http.ResponseWriter, r *http.Request) {
	// 从URL中获得uuid
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 根据uuid从磁盘读取tempInfo
	tempInfo, err := getTempInfo(uuid)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.Remove(os.Getenv("STORAGE_ROOT") + "/temp/" + tempInfo.Uuid)
	// 检查临时文件的大小是否与请求头中的size相同
	stat, _ := os.Stat(os.Getenv("STORAGE_ROOT") + "/temp/" + tempInfo.Uuid + ".dat")
	fmt.Println(stat.Size(), tempInfo.Size)
	if stat.Size() != tempInfo.Size {
		log.Println("size not match")
		os.Remove(os.Getenv("STORAGE_ROOT") + "/temp/" + tempInfo.Uuid + ".dat")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 对临时文件基于gzip进行压缩
	readFile, _ := os.OpenFile(os.Getenv("STORAGE_ROOT")+"/temp/"+tempInfo.Uuid+".dat", os.O_CREATE|os.O_APPEND, 6) 
	defer readFile.Close()
	buf,_:=ioutil.ReadAll(readFile)
	log.Println(string(buf))
	utils.Gzip(buf, string(os.Getenv("STORAGE_ROOT")+"/objects/"+tempInfo.Hash+"."+strconv.Itoa(tempInfo.ID)))
	os.Remove(os.Getenv("STORAGE_ROOT")+"/temp/"+tempInfo.Uuid+".dat")

	log.Println("add object:  "+ tempInfo.Hash + "." + strconv.Itoa(tempInfo.ID)+" in "+ string(os.Getenv("STORAGE_ROOT")+"/objects/"+tempInfo.Hash+"."+strconv.Itoa(tempInfo.ID)))
	locate.Add(tempInfo.Hash, tempInfo.ID)
}

func delete(w http.ResponseWriter, r *http.Request) {
	// 从URL中获得uuid
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]

	// 将保存tempInfo的文件和临时文件删除
	os.Remove(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid + ".dat")
	os.Remove(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid)
}

func get(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	f, e := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid + ".dat")
	defer f.Close()
	if e != nil {
		log.Println(e.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(w, f)
}

func head(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	stat, e := os.Stat(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid + ".dat")
	if e != nil {
		log.Println(e.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// log.Println("current size:",stat.Size())
	w.Header().Set("length", strconv.FormatInt(stat.Size(), 10))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		post(w, r)
	case http.MethodPatch:
		patch(w, r)
	case http.MethodPut:
		put(w, r)
	case http.MethodDelete:
		delete(w, r)
	case http.MethodGet:
		get(w, r)
	case http.MethodHead:
		head(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
