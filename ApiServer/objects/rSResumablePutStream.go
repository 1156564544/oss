package objects

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"httpTool"
	"log"
	"net/http"

	"rs"
)

type PutToken struct {
	Name string
	Hash string
	Size int64
	Servers []string
	Uuids []string
}

type RSResumablePutStream struct {
	RSPutStream *RSPutStream
	Token *PutToken
}

func CreateRSResumablePutStream(dataServers []string, name string,hash string, size int64) (*RSResumablePutStream, error) {
	// 首先创造RSPutStream
	rsPutStream,err:=CreateRSPutStream(dataServers,hash,size)
	if err!=nil{
		return nil,err
	}
	// 然后把各个DataServer的uuid读出来
	uuids:=make([]string,rs.NUM_SHARDS)
	for i:=range uuids{
		uuids[i]=rsPutStream.writers[i].uuid
	}
	// 最后生成token
	token:=&PutToken{name,hash,size,dataServers,uuids}
	return &RSResumablePutStream{rsPutStream,token},nil
}

func GetRSResumablePutStreamFromToken(token string) (*RSResumablePutStream,error) {
	decoded,_:=base64.StdEncoding.DecodeString(token)
	var putToken PutToken
	err:=json.Unmarshal(decoded,&putToken)
	if err!=nil{
		return nil,err
	}
	writers:=make([]*putStream,rs.NUM_SHARDS)
	for i:=range writers{
		writers[i]=&putStream{putToken.Uuids[i],putToken.Servers[i]}
	}
	rsPutStream,err:=CreateRSPutStreamFromPutStreams(writers)
	if err!=nil{
		return nil,err
	}
	return &RSResumablePutStream{rsPutStream,&putToken},nil
}

func (stream *RSResumablePutStream) ToToken() string {
	token,_:=json.Marshal(stream.Token)
	return base64.StdEncoding.EncodeToString(token)
}

func (s *RSResumablePutStream) CurrentSize() int64 {
	r, e := http.Head(fmt.Sprintf("http://%s/temp/%s", s.Token.Servers[0], s.Token.Uuids[0]))
	if e != nil {
		log.Println(e)
		return -1
	}
	if r.StatusCode != http.StatusOK {
		log.Println(r.StatusCode)
		return -1
	}
	size := httpTool.GetSizeFromHeader(r.Header) * rs.NUM_DATA_SHARES
	if size > s.Token.Size {
		size = s.Token.Size
	}
	return size
}