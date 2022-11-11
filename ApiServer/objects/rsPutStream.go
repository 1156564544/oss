package objects

import (
	"fmt"
	"sync"
	"errors"
	reedsolomon "github.com/klauspost/reedsolomon"

	"rs"
)

type RSPutStream struct {
	writers []*putStream
	enc    reedsolomon.Encoder
	cache []byte
}

func CreateRSPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	if len(dataServers) < rs.NUM_DATA_SHARES+rs.NUM_PARITY_SHARES {
		return nil, errors.New("no enough dataServers")
	}
	writers := make([]*putStream, rs.NUM_DATA_SHARES+rs.NUM_PARITY_SHARES)
	perShard:=getPersharedSize(size)
	var wg sync.WaitGroup
	wg.Add(rs.NUM_DATA_SHARES+rs.NUM_PARITY_SHARES)
	for i:=0;i<rs.NUM_DATA_SHARES+rs.NUM_PARITY_SHARES;i++{
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
			err=stream.Flush()
			if err!=nil{
				return 0,err
			}
			break
		}
		need:=rs.CHUNK_SIZE*rs.NUM_DATA_SHARES-len(stream.cache)
		stream.cache=append(stream.cache,p[current:current+need]...)
		stream.Flush()
		if err!=nil{
			return 0,err
		}
		current+=need
	}
	return len(p),nil
}

func (stream *RSPutStream) Flush() error {
	shards,_:=stream.enc.Split(stream.cache)
	stream.enc.Encode(shards)
	for i := 0; i < rs.NUM_DATA_SHARES+rs.NUM_PARITY_SHARES; i++ {
		_,err:=stream.writers[i].Write(shards[i])
		if err!=nil{
			return fmt.Errorf("%v-th putStream write error: %v", i, err.Error())
		}
	}
	stream.cache=stream.cache[:0]
	return nil
}

func (stream *RSPutStream)commit(iscommit bool) error{
	var wg sync.WaitGroup
	wg.Add(rs.NUM_DATA_SHARES+rs.NUM_PARITY_SHARES)
	for i:=0;i<rs.NUM_DATA_SHARES+rs.NUM_PARITY_SHARES;i++{
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
