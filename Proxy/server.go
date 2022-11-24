package main

import (
	"os"
	"log"
	"net/http"

	"Proxy/login"
	"Proxy/manage"
)

func main(){
	ip := os.Args[1]
	http.HandleFunc("/login", login.Handler)
	http.HandleFunc("/addUser", manage.AddUser)
	http.HandleFunc("/deleteUser", manage.DeleteUser)
	http.HandleFunc("/updateUser", manage.UpdateUser)
	log.Fatal(http.ListenAndServeTLS(ip, "./secret/ca.crt", "./secret/ca.key", nil))

}