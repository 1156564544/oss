package objects

import (
	"errors"
	"fmt"
	"log"
	"sync"

	reedsolomon "github.com/klauspost/reedsolomon"

	"rs"
)

type RSPutStream struct {
	writers []*putStream
	enc    reedsolomon.Encoder
	cache []byte
}

func CreateRSPutStreamFromPutStreams(writers []*putStream) (*RSPutStream,error){
	if len(writers)!=rs.NUM_SHARDS{
		return nil,errors.New("Invalid writers")
	}
	enc, err := reedsolomon.New(rs.NUM_DATA_SHARES, rs.NUM_PARITY_SHARES)
	if err != nil {
		return nil, err
	}
	return &RSPutStream{writers,enc,make([]byte,0)},nil
}

func CreateRSPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	if len(dataServers) < rs.NUM_SHARDS {
		return nil, errors.New("no enough dataServers")
	}
	writers := make([]*putStream, rs.NUM_SHARDS)
	perShard:=getPersharedSize(size)
	var wg sync.WaitGroup
	wg.Add(rs.NUM_SHARDS)
	for i:=0;i<rs.NUM_SHARDS;i++{
		go func (i int){
			var err error
			writers[i],err=CreatePutStream(dataServers[i],fmt.Sprintf("%s.%v",hash,i),perShard)
			if err!=nil{
				panic(err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	enc,_:=reedsolomon.New(rs.NUM_DATA_SHARES,rs.NUM_PARITY_SHARES)
	return &RSPutStream{writers,enc,make([]byte,0)},nil
}

func (stream *RSPutStream) Write(p []byte) (n int, err error) {
	current:=0
	for current<len(p){
		if len(stream.cache)+len(p[current:])<=rs.CHUNK_SIZE*rs.NUM_DATA_SHARES{
			stream.cache=append(stream.cache,p[current:]...)
			log.Println("flush1")
			err=stream.Flush()
			if err!=nil{
				return 0,err
			}
			break
		}
		need:=rs.CHUNK_SIZE*rs.NUM_DATA_SHARES-len(stream.cache)
		stream.cache=append(stream.cache,p[current:current+need]...)
		// stream.Flush()
		if len(stream.cache)==rs.CHUNK_SIZE*rs.NUM_DATA_SHARES{
			log.Println("flush2")
			stream.Flush()
		}
		if err!=nil{
			return 0,err
		}
		current+=need
	}
	return len(p),nil
}

func (stream *RSPutStream) Flush() error {
	shards,_:=stream.enc.Split(stream.cache)
	log.Println("stream.cache:",stream.cache)
	log.Println("shards:",shards)
	stream.enc.Encode(shards)
	for i := 0; i < rs.NUM_SHARDS; i++ {
		_,err:=stream.writers[i].Write(shards[i])
		if err!=nil{
			return fmt.Errorf("%v-th putStream write error: %v", i, err.Error())
		}
	}
	stream.cache=stream.cache[:0]
	return nil
}

func (stream *RSPutStream)Commit(iscommit bool) error{
	var wg sync.WaitGroup
	wg.Add(rs.NUM_SHARDS)
	for i:=0;i<rs.NUM_SHARDS;i++{
		go func (i int){
			err:=stream.writers[i].commit(iscommit)
			if err!=nil{
				panic(err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return nil
}
