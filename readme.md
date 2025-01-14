package main

import (
"fmt"
ninaRouter "github.com/JonecoBoy/nina/router"
"net/http"
"time"
)

type User struct {
Username string `json:"username"`
Password string `json:"password"`
}

var dummyUser = User{
Username: "admin",
Password: "password",
}

func main() {
// todo eventos
// todo presentation para dar match com o return application. e add parser
// add viper pro .env
// todo criar roteador  / grupo no router tpo /joneco sub routes
roteador := ninaRouter.NewRouter()

	roteador.GET("/hello/{id}/{abc}", helloHandler, nil)

	adminGroup := roteador.GROUP("/admin", nil, nil)

	adminGroup.GET("/", func(w http.ResponseWriter, r *ninaRouter.NinaRequest) { fmt.Println("Request to / from group") }, nil)
	adminGroup.GET("/hello", helloHandler, nil)
	adminGroup.POST("/hello", func(w http.ResponseWriter, r *ninaRouter.NinaRequest) {
		fmt.Fprint(w, "Hello from post")
	}, []ninaRouter.Middleware{})

	fmt.Println("Server is running on port 8081")
	http.ListenAndServe(":8081", roteador)

}

func helloHandler(w http.ResponseWriter, r *ninaRouter.NinaRequest) {
fmt.Fprint(w, "Hello from get")
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
cookie := http.Cookie{
Name:     "exampleCookie",
Value:    "cookieValue",
Expires:  time.Now().Add(24 * time.Hour),
HttpOnly: true,
}
http.SetCookie(w, &cookie)
w.Write([]byte("Cookie set!"))
}
