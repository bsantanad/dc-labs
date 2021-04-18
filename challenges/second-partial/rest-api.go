package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type Image struct {
	Name string `json:"name"`
	Size int    `json:"size"`
	Data []byte `json:"data"`
}

type User struct {
	Username string  `json:"user"`
	Token    string  `json:"token"`
	Images   []Image `json:"image"`
	Time     string  `json:"time"`
}

type Status struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

type ImageMsg struct {
	Message  string `json:"message"`
	Filename string `json:"filename"`
	Size     string `json:"size"`
}

type Message struct {
	Message string `json:"message"`
}

var Users []User /* this will act as our DB */

/********************* Endpoint Functions ***************************/

func homePage(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(200)
	returnMsg(w, "DPIP REST API index. Invalid enpoints will redirect here")
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
		Time:     time.Now().UTC().String(),
	}
	Users = append(Users, userInfo)

	json.NewEncoder(w).Encode(login)
}

// delLogout function will revoke a token from being usable.
// first it checks if the headers are sent in the correct
// format, then it will search the token in the Users "DB"
// if found it will remove it, if not, it will return 400
func delLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[INFO]: DELETE /logout requested")
	tmp := r.Header.Get("Authorization")
	if strings.Fields(tmp)[0] != "Bearer" {
		w.WriteHeader(400)
		returnMsg(w, "bad request, check headers "+
			"you must send a Bearer token")
		return
	}
	token := strings.Fields(tmp)[1] // get the token from header
	index, user, exists := searchToken(token)
	if !exists {
		w.WriteHeader(400)
		returnMsg(w, "token not found, "+
			"please provide a valid one")
		return
	}

	Users = removeUser(Users, index)
	returnMsg(w, "Bye "+user.Username+", your token has been revoked")
}

// based on https://stackoverflow.com/a/40699578
// postUpload, upload a file (image).
// It first checks the headers and find the token,
// validates it and finds the user.
// Then creates a buffer, copy the bytes of the image
// to it and fills the Image struct.
// Finally it append the image to the Image slice
// the user has.
func postUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[INFO]: POST /upload requested")
	tmp := r.Header.Get("Authorization")
	if strings.Fields(tmp)[0] != "Bearer" {
		w.WriteHeader(400)
		returnMsg(w, "bad request, check headers "+
			"you must send a Bearer token")
		return
	}
	token := strings.Fields(tmp)[1] // get the token from header
	index, user, exists := searchToken(token)
	if !exists {
		w.WriteHeader(400)
		returnMsg(w, "token not found, "+
			"please provide a valid one")
		return
	}

	// uploading the file part
	r.ParseMultipartForm(32 << 20) // limit your max input length!
	var buf bytes.Buffer
	file, header, err := r.FormFile("data")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	// Copy the image data to my buffer
	io.Copy(&buf, file)

	// Fill the image struct
	var image Image
	image.Name = name[0]
	image.Size = buf.Len()
	image.Data, err = buf.ReadBytes(254)
	if err != nil {
		w.WriteHeader(409)
		returnMsg(w, "Image couldn't be uploaded :(. Please try again")
		return
	}
	Users[index].Images = append(user.Images, image)

	buf.Reset()

	var msg ImageMsg
	msg = ImageMsg{
		Message:  "An image has been successfully uploaded :)",
		Filename: image.Name,
		Size:     fmt.Sprintf("%d bytes", image.Size),
	}

	json.NewEncoder(w).Encode(msg)
}

// getStatus, show the status of the account related
// to the token sent in the header, proper validations
// are done, and then the creation time, and a msg is
// returned to the user
func getStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[INFO]: GET /status requested")
	tmp := r.Header.Get("Authorization")
	if strings.Fields(tmp)[0] != "Bearer" {
		w.WriteHeader(400)
		returnMsg(w, "bad request, check headers "+
			"you must send a Bearer token")
		return
	}
	token := strings.Fields(tmp)[1] // get the token from header
	_, user, exists := searchToken(token)
	if !exists {
		w.WriteHeader(400)
		returnMsg(w, "token not found, "+
			"please provide a valid one")
		return
	}

	var status Status
	status = Status{
		Message: "Hi " + user.Username + ", the DPIP System is Up and Running",
		Time:    user.Time,
	}

	json.NewEncoder(w).Encode(status)
}

/********************* Handler Functions ***************************/

func handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
	case http.MethodPost:
		postLogin(w, r) // post
	case http.MethodPut:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
	case http.MethodDelete:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
	default:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
	}

}
func handleLogout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
	case http.MethodPost:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
	case http.MethodPut:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
	case http.MethodDelete:
		delLogout(w, r) // delete
	default:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
	}

}
func handleUpload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
	case http.MethodPost:
		postUpload(w, r) // post
	case http.MethodPut:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
	case http.MethodDelete:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
	default:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
	}

}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getStatus(w, r) //get
	case http.MethodPost:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
	case http.MethodPut:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
	case http.MethodDelete:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
	default:
		w.WriteHeader(404)
		returnMsg(w, "page not found")
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

/********************* Helper Functions ***************************/

// Search token in Users, returned index, user struct
// and boolean that tells us if it was found.
func searchToken(token string) (int, User, bool) {
	for i, user := range Users {
		if user.Token == token {
			return i, user, true
		}
	}
	var tmp User
	return -1, tmp, false
}

// swap the user you want to remove with the
// last item, return the slice without the last item
func removeUser(users []User, index int) []User {
	users[index] = users[len(users)-1]
	return users[:len(users)-1]
}

func returnMsg(w http.ResponseWriter, msg string) {
	var msgJSON Message
	msgJSON = Message{
		Message: msg,
	}
	json.NewEncoder(w).Encode(msgJSON)

}

func main() {
	handleRequests()
}
