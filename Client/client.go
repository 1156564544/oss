package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"net/url"
)

func calculateHash(s string) string {
	r := strings.NewReader(s)
	h := sha256.New()
	io.Copy(h, r)
	hash := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return url.PathEscape(hash)
}

func PostHeader(url string, msg string, headers map[string]string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(""))
	if err != nil {
		return "", err
	}
	fmt.Println(headers)
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	token := resp.Header.Get("location")
	return token, nil
}

func post(object string, msg string) (string, error) {
	hash := calculateHash(msg)
	url := "http://localhost:10000/objects/" + object
	headers := make(map[string]string)
	headers["Digest"] = "SHA-256=" + hash
	headers["length"] = strconv.Itoa(len(msg))
	token, err := PostHeader(url, msg, headers)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return token, nil
}

func putFromToken(token string, offset int64, length int64,msg string) error {
	client := &http.Client{}

	url := "http://localhost:10000" + token
	req, err := http.NewRequest("PUT", url, strings.NewReader(msg[offset:offset+length]))
	if err != nil {
		return err
	}
	headers := make(map[string]string)
	headers["Range"] = "bytes=" + strconv.FormatInt(offset, 10) + "_" + strconv.FormatInt(offset+length-1, 10) + "/" + strconv.FormatInt(int64(len(msg)), 10)
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func put(object string, msg string) error {
	client := &http.Client{}

	url := "http://localhost:10000/objects/" + object
	req, err := http.NewRequest("PUT", url, strings.NewReader(msg))
	if err != nil {
		return err
	}
	headers := make(map[string]string)
	headers["length"] = strconv.Itoa(len(msg))
	headers["Digest"] = "SHA-256=" + calculateHash(msg)
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func get(object string) string {
	url := "http://localhost:10000/objects/" + object
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	buf := make([]byte, 1024)
	n, err := resp.Body.Read(buf)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(buf[:n])
}

func main() {
	object := "test6_24"
	msg := "this is the object test6_24"
	// token, err := post(object, msg)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }
	// fmt.Println(token)
	// err = putFromToken(token, 0, 8,msg)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }
	// err = putFromToken(token, 8, int64(len(msg))-8,msg)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }
	put(object, msg)
	// get(object)
}
