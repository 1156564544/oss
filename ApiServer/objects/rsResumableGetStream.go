package objects

import (
	"errors"
	reedsolomon "github.com/klauspost/reedsolomon"
	"io"
	"log"
	"net/http"
	"strconv"

	"rs"
)

func CreateResumableGetStream(server, name string, size int64) (*getStream, error) {
	resp, err := http.Get("http://localhost" + server + "/temp/" + name)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("resp.StatusCode:", resp.StatusCode)
		return nil, errors.New("Get with status code " + strconv.Itoa(resp.StatusCode))
	}
	return &getStream{resp.Body}, nil
}

func CreateRSResumableGetStream(dataServers []string, uuids []string, size int64) (*RSGetStream, error) {
	writers := make([]io.Writer, rs.NUM_SHARDS)
	readers := make([]io.Reader, rs.NUM_SHARDS)
	for i := range readers {
		var e error
		readers[i], e = CreateResumableGetStream(dataServers[i], uuids[i], size)
		if e != nil {
			return nil, e
		}
	}
	enc, _ := reedsolomon.New(rs.NUM_DATA_SHARES, rs.NUM_PARITY_SHARES)
	return &RSGetStream{writers: writers, readers: readers, enc: enc, size: size, cache: make([]byte, 0), total: 0}, nil
}
