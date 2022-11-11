package main

import (
	"DataServer/heartbeat"
	"DataServer/locate"
	"DataServer/objects"
	"DataServer/temp"
	"log"
	"net/http"
	"os"
)

func main() {
	ip := os.Args[1]
	log.Printf("DataServer_%v start...\n", ip)
	go heartbeat.StartHeartbeat(ip)
	go locate.StartLocate(ip)
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(ip, nil))
}