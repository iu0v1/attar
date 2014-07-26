Attar
=====

Go pakage for http session auth.

Pkg for use with gorilla/mux router.

Get Attar
=====
```
go get github.com/SpiritOfStallman/attar
```

Usage
=====

```Go
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
```

### User AuthProvider
User functon must take 'user' and 'password' arguments, and return true (if user auth successfully) or false (if auth data false). 
```Go
// user code
func checkAuth(u, p string) bool {
    if u == "user" && p == "qwerty" {
        return true
    }
    return false
}

// and define it
a := attar.New()
a.SetAuthProvider(checkAuth)
```
Also attar include pre-define simple AuthProvider:
```Go
// users list based on map[user]password
userList := map[string]string{
    "user":  "qwerty",
    "admin": "asdfgh",
}

a := attar.New()
a.SetAuthProvider(a.SimpleAuthProvider(userList))
```

### User cookie keys
Attar can create new sessions by keys:
```Go
import (
    "github.com/SpiritOfStallman/attar"
)

func main() {
    ...
    a := attar.New()
    a.SetCookieSessionKeys(
        []byte("261AD9502C583BDQQQQQQQQQQQQQQQQQ"),
        []byte("RRRRRRRRRRRRRRR3FC5C7B3D6E4DDAFF"),
    )
    ...
}
```

And can use existing 'gorilla/sessions' CookieStore:
```Go
import (
    "github.com/gorilla/sessions"
    "github.com/SpiritOfStallman/attar"
)

func main() {
    ..
    gorillaSessions := sessions.NewCookieStore(
        []byte("261AD9502C583BD7D8AA03083598653B"),
        []byte("E9F6FDFAC2772D33FC5C7B3D6E4DDAFF"),
    )
    ..
    a := attar.New()
    a.SetGorillaCookieStore(gorillaSessions)
    ..
}
```

DOC
=====
For more information refer to pkg doc http://godoc.org/github.com/SpiritOfStallman/attar.
