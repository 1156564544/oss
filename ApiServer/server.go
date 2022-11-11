package main

import (
	"ApiServer/heartbeat"
	"ApiServer/locate"
	"ApiServer/objects"
	"ApiServer/versions"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ip := os.Args[1]
	go heartbeat.ListenHeartbeat()
	go func(){
		time.Sleep(5*time.Second)
		log.Printf("ApiServer_%v start...\n", ip)
	}()
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(ip, nil))
}
