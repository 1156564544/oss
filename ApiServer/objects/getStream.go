package objects

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

type getStream struct {
	reader io.Reader
}

func CreateGetStream(server, name string, size int64) (*getStream, error) {
	resp, err := http.Get("http://localhost" + server + "/objects/" + name)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("resp.StatusCode:", resp.StatusCode)
		return nil, errors.New("Get with status code " + strconv.Itoa(resp.StatusCode))
	}
	return &getStream{resp.Body},nil
}

func (r *getStream) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}