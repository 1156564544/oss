package main

import (
	"DataServer/heartbeat"
	"DataServer/locate"
	"DataServer/objects"
	"log"
	"net/http"
	"os"
)

func locate2(name string) bool {
	_, err := os.Stat(name)
	return err == nil || !os.IsNotExist(err)
}

func main() {
	ip := os.Args[1]
	go heartbeat.StartHeartbeat(ip)
	go locate.StartLocate(ip)
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(ip, nil))
}
