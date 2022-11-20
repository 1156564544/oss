package utils

import (
	"compress/gzip"
	"io"
	"os"
	"net/http"
)
// 压缩
func Gzip(content []byte, path string) error {
	gzFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer gzFile.Close()
	gzipWriter := gzip.NewWriter(gzFile)

	defer gzipWriter.Close()
	//gzipWriter.Name = fileName
	_, err = gzipWriter.Write(content)
	if err != nil {
		return err
	}
	return nil
}

// 解压
func UnGzip(path string,w *http.ResponseWriter) (err error) {
	gzipFile, err := os.Open(path)
	if err != nil {
		return 
	}
	defer gzipFile.Close()
	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return 
	}
	defer gzipReader.Close()
	_,err=io.Copy(*w, gzipReader)
	if err != nil {
		return
	}
	// var buf bytes.Buffer
	//_, err = io.Copy(&buf, gzipReader)
	//if err != nil {
	//	return err
	//}
	return nil
}