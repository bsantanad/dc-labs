package main

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type User struct {
	Username string `json:"user"`
	Token    string `json:"token"`
}

var Users []User /* this will act as out DB */

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "distributed and parallel image processing rest api\n")
	fmt.Println("[INFO]: / requested")
}

// postLogin will get the hash that's generated by default
// by the header "Authorization", then it will use it
// as the token for this particular user.
// It will also add the user to the "DB" of users, along
// with it's token
func postLogin(w http.ResponseWriter, r *http.Request) {
	var token string
	var user string
	var tmp string

	fmt.Println("[INFO]: POST /login requested")
	user, _, _ = r.BasicAuth() //get username
	tmp = r.Header.Get("Authorization")
	token = strings.Fields(tmp)[1] // get the hash from header

	//Build response
	var login LoginResponse
	login = LoginResponse{
		Message: "Hi " + user + ", welcome to the DPIP System",
		Token:   token,
	}

    var userInfo User
	userInfo = User{
		Username: user,
		Token:    token,
	}
	Users = append(Users, userInfo)

	json.NewEncoder(w).Encode(login)
}

func delLogout(w http.ResponseWriter, r *http.Request) {
	return
}
func postUpload(w http.ResponseWriter, r *http.Request) {
	return
}
func getStatus(w http.ResponseWriter, r *http.Request) {
	return
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.Error(w, "not found", 404)
	case http.MethodPost:
		postLogin(w, r) // post
	case http.MethodPut:
		http.Error(w, "not found", 404)
	case http.MethodDelete:
		http.Error(w, "not found", 404)
	default:
		http.Error(w, "not found", 404)
	}

}
func handleLogout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.Error(w, "not found", 404)
	case http.MethodPost:
		http.Error(w, "not found", 404)
	case http.MethodPut:
		http.Error(w, "not found", 404)
	case http.MethodDelete:
		delLogout(w, r) // delete
	default:
		http.Error(w, "not found", 404)
	}

}
func handleUpload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.Error(w, "not found", 404)
	case http.MethodPost:
		postUpload(w, r) // post
	case http.MethodPut:
		http.Error(w, "not found", 404)
	case http.MethodDelete:
		http.Error(w, "not found", 404)
	default:
		http.Error(w, "not found", 404)
	}

}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getStatus(w, r) //get
	case http.MethodPost:
		http.Error(w, "not found", 404)
	case http.MethodPut:
		http.Error(w, "not found", 404)
	case http.MethodDelete:
		http.Error(w, "not found", 404)
	default:
		http.Error(w, "not found", 404)
	}

}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/status", handleStatus)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
