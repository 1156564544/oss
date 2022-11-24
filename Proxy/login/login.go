package login

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"log"
	"net/http"

	"Proxy/utils"
	"redisTool"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	log.Println("username:", username, "password:", password)
	if ok {
		// 根据用户名从数据库中获得用户信息
		user, err := utils.SelectFromDB(username)
		if err != nil {
			log.Println(err)
			w.Write([]byte("no such user"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		expectedPassword:=user.Password
		// Calculate SHA-256 hashes for the provided and expected
		// usernames and passwords.
		passwordHash := sha256.Sum256([]byte(password))
		expectedPasswordHash := sha256.Sum256([]byte(expectedPassword))

		// 使用 subtle.ConstantTimeCompare() 进行校验
		// the provided username and password hashes equal the  
		// expected username and password hashes. ConstantTimeCompare
		// 如果值相等，则返回1，否则返回0。
		// Importantly, we should to do the work to evaluate both the 
		// username and password before checking the return values to 
		// 避免泄露信息。
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

		// If the username and password are correct, then call
		// the next handler in the chain. Make sure to return 
		// afterwards, so that none of the code below is run.
		if !passwordMatch {
			log.Println("password not match")
			w.Write([]byte("password is wrong"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		b,_:=json.Marshal(*user)
		token,err:=utils.EncryptByAes([]byte(b))
		if err!=nil{	
			log.Println("EncryptByAes fail")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Authorization",token)
		w.WriteHeader(http.StatusOK)
		redisTool.SetAdd(token)
		return
	}

	// If the Authentication header is not present, is invalid, or the
	// username or password is wrong, then set a WWW-Authenticate 
	// header to inform the client that we expect them to use basic
	// authentication and send a 401 Unauthorized response.
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}