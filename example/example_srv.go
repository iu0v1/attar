package main

import (
	"html/template"
	"net/http"

	"github.com/SpiritOfStallman/attar"
	"github.com/gorilla/mux"
)

// main page
var mainPage = template.Must(template.New("").Parse(`
	<html><head></head><body><center>
	<h1 style="padding-top:15%;">HELLO!</h1>
	</form></center></body>
	</html>`))

func mainPageHandler(res http.ResponseWriter, req *http.Request) {
	mainPage.Execute(res, nil)
}

// login page
var loginPage = template.Must(template.New("").Parse(`
	<html><head></head><body>
	<center>
	<form id="login_form" action="/login" method="POST" style="padding-top:15%;">
	<p>user::qwerty</p>
	<input type="text" name="login" placeholder="Login" autofocus><br>
	<input type="password" placeholder="Password" name="password"><br>
	<input type="submit" value="LOGIN">
	</form></center></body>
	</html>`))

func loginPageHandler(res http.ResponseWriter, req *http.Request) {
	loginPage.Execute(res, nil)
}

// auth provider function
func checkAuth(u, p string) bool {
	if u == "user" && p == "qwerty" {
		return true
	}
	return false
}

func main() {

	a := attar.New()

	a.SetAuthProvider(checkAuth)
	a.SetLoginRoute("/login")
	a.SetCookieSessionKeys(
		[]byte("261AD9502C583BDQQQQQQQQQQQQQQQQQ"),
		[]byte("RRRRRRRRRRRRRRR3FC5C7B3D6E4DDAFF"),
	)

	// set options, with session & cookie lifetime == 30 sec
	options := &attar.AttarOptions{
		Path:                       "/",
		MaxAge:                     60,
		HttpOnly:                   true,
		SessionName:                "test-session",
		SessionLifeTime:            60,
		SessionBindUseragent:       true,
		SessionBindUserHost:        true,
		LoginFormUserFieldName:     "login",
		LoginFormPasswordFieldName: "password",
	}
	a.SetAttarOptions(options)

	// create mux router
	router := mux.NewRouter()
	router.HandleFunc("/", mainPageHandler)
	router.HandleFunc("/login", loginPageHandler).Methods("GET")
	// set attar.AuthHandler as handler func
	// for check login POST data
	router.HandleFunc("/login", a.AuthHandler).Methods("POST")

	// set auth proxy function
	http.Handle("/", a.GlobalAuthProxy(router))

	// start net/httm server at 8080 port
	if err := http.ListenAndServe("127.0.0.1:8082", nil); err != nil {
		panic(err)
	}
}
