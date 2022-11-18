package objects

import (
	"errors"
	"fmt"
	"io"

	reedsolomon "github.com/klauspost/reedsolomon"

	"ApiServer/heartbeat"
	"ApiServer/locate"
	"rs"
)

type RSGetStream struct {
	writers []io.Writer
	readers []io.Reader
	enc     reedsolomon.Encoder
	size    int64
	cache   []byte
	total   int64
}

func CreateRSGetStream(hash string, size int64) (*RSGetStream, error) {
	ips := locate.Locate(hash)
	// chunk数目少于编码时数据块个数，无法进行纠删码
	if len(ips) < rs.NUM_DATA_SHARES {
		return nil, errors.New("Not enough dataServers")
	}
	dataServers := make([]string, 0)
	// 有chunk丢失了，需要恢复丢失的chunk
	if len(ips) < rs.NUM_SHARDS {
		dataServers = heartbeat.RandomChooseDataServersWithExclude(rs.NUM_SHARDS, ips)
	}
	writers := make([]io.Writer, rs.NUM_SHARDS)
	readers := make([]io.Reader, rs.NUM_SHARDS)
	perShard := getPersharedSize(size)
	for i := 0; i < rs.NUM_SHARDS; i++ {
		if server, ok := ips[i]; !ok {
			server = dataServers[0]
			dataServers = dataServers[1:]
			var err error
			writers[i], err = CreatePutStream(server, fmt.Sprintf("%s.%v", hash, i), perShard)
			if err != nil {
				return nil, err
			}
		} else {
			var err error
			readers[i], err = CreateGetStream(server, fmt.Sprintf("%s.%v", hash, i), size)
			if err != nil {
				return nil, err
			}
		}
	}
	enc, _ := reedsolomon.New(rs.NUM_DATA_SHARES, rs.NUM_PARITY_SHARES)
	return &RSGetStream{writers: writers, readers: readers, enc: enc, size: size, cache: make([]byte, 0), total: 0}, nil
}

func (stream *RSGetStream) Read(p []byte) (int, error) {
	// 当前cache已空，需要重新从dataServers拉取数据来解码填充到cache
	if len(stream.cache) == 0 {
		err := stream.decoding()
		if err != nil {
			return 0, err
		}
	}
	// 从cache中读取数据
	length := len(p)
	if length > len(stream.cache) {
		length = len(stream.cache)
	}
	copy(p, stream.cache[:length])
	stream.cache = stream.cache[length:]
	return length, nil
}

func (stream *RSGetStream) decoding() error {
	// object已经读完
	if stream.total >= stream.size {
		return io.EOF
	}
	// shards: 读取数据块
	shards := make([][]byte, rs.NUM_SHARDS)
	// repareIds: 记录需要恢复的chunk号
	repareIds := make([]int, 0)
	// 从dataServers读取数据
	for i := range shards {
		if stream.readers[i] != nil {
			shards[i] = make([]byte, rs.CHUNK_SIZE)
			n, e := io.ReadFull(stream.readers[i], shards[i])
			if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
				shards[i] = nil
			} else if n != rs.CHUNK_SIZE {
				shards[i] = shards[i][:n]
			}
		} else {
			repareIds = append(repareIds, i)
		}
	}
	// 根据shareds重构原数据
	err := stream.enc.Reconstruct(shards)
	// log.Println("After repair:", shards)
	if err != nil {
		return fmt.Errorf("Reconstruct error: %v", err.Error())
	}
	// 将恢复的数据写入dataServers
	for _, id := range repareIds {
		stream.writers[id].Write(shards[id])
	}
	// 将数据写入cache以供read读取
	for i := 0; i < rs.NUM_DATA_SHARES; i++ {
		length := int64(len(shards[i]))
		if stream.total+length > stream.size {
			length = stream.size - stream.total
		}
		stream.cache = append(stream.cache, shards[i][:length]...)
		stream.total += length
	}
	return nil
}

func (stream *RSGetStream) Close() error {
	for _,writer:=range stream.writers{
		if writer!=nil{
			writer.(*putStream).commit(true)
		}
	}
	return nil
}

func (stream *RSGetStream) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekCurrent {
		return 0, errors.New("Only support io.SeekCurrent")
	}
	if offset < 0 {
		return 0, errors.New("Offset must be positive")
	}
	for offset>0{
		length:=getRoundSize()
		if length>int(offset){
			length=int(offset)
		}
		buf:=make([]byte,length)
		io.ReadFull(stream, buf)
		offset-=int64(length)
	}
	return offset, nil 
}