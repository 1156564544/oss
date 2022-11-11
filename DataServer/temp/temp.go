package temp

import (
	"DataServer/locate"
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
	ID  int
	Size int64
}


// 处理来自apiServer的post请求，其中包括待上传的对象名（<hash>.<id of shard>)和单个shard的大小，dataServer生成uuid并作为post请求的响应返回给apiServer。
// dataServer将uuid作为文件名保存在temp目录下，文件内容为json格式的tempInfo，其中包括uuid、hash、id和size。
func post(w http.ResponseWriter, r *http.Request) {
	object_name:= strings.Split(r.URL.EscapedPath(), "/")[2]
	fmt.Println(object_name)
	// object_name: <hash>.<id of shard>
	hash:=strings.Split(object_name, ".")[0]
	id,_:= strconv.Atoi(strings.Split(object_name, ".")[1])
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
		ID:  id,
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
	f,err:=os.OpenFile(os.Getenv("STORAGE_ROOT")+"/temp/"+tempInfo.Uuid+".dat", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	defer f.Close()
	if err != nil {
		log.Println("Create file failed: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
	stat,_:=os.Stat(os.Getenv("STORAGE_ROOT")+"/temp/"+tempInfo.Uuid+".dat")
	fmt.Println(stat.Size(),tempInfo.Size)
	if stat.Size()!=tempInfo.Size{
		log.Println("size not match")
		os.Remove(os.Getenv("STORAGE_ROOT")+"/temp/"+tempInfo.Uuid+".dat")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 将临时文件重命名为正式文件
	err = os.Rename(os.Getenv("STORAGE_ROOT")+"/temp/"+tempInfo.Uuid+".dat", string(os.Getenv("STORAGE_ROOT")+"/objects/"+tempInfo.Hash+"."+strconv.Itoa(tempInfo.ID)))
	// err=os.Rename(os.Getenv("STORAGE_ROOT")+"/temp/"+tempInfo.Uuid+".dat", "/data1/objects/ic%2F8+5zNfK4R9PY5wlfgAGl8wWF8aEORabXRGReyMXg=")
	if err != nil {
		log.Println("Rename failed: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	locate.Add(tempInfo.Hash,tempInfo.ID)
}

func delete(w http.ResponseWriter, r *http.Request) {
	// 从URL中获得uuid
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]

	// 将保存tempInfo的文件和临时文件删除
	os.Remove(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid + ".dat")
	os.Remove(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid)
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}