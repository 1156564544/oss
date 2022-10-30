package main

import (
	"ApiServer/heartbeat"
	"ApiServer/locate"
	"ApiServer/objects"
	"ApiServer/versions"
	"log"
	"net/http"
	"os"
)

func main() {
	ip := os.Args[1]
	go heartbeat.ListenHeartbeat()
	
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(ip, nil))
}
