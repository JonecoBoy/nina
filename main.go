package main

import (
	"fmt"
	"net/http"
	ninaMiddleware "nina/middleware"
	ninaRouter "nina/router"
)

func main() {
	// todo eventos
	// todo presentation para dar match com o return application. e add parser
	// add viper pro .env
	// todo criar roteador  / grupo no router tpo /joneco sub routes
	roteador := ninaRouter.NewRouter()

	roteador.GET("/hello/{id}/{abc}", heloHandler, ninaMiddleware.LogginMiddleware)

	//roteador.POST("/hello/{id}", heloHandler)
	//router.Handle("/hello/{version}", middleware(router.HandlerFunc(heloHandler)))

	//router.HandleFunc("/", func(writer router.ResponseWriter, request *router.Request) {
	//	fmt.Fprint(writer, "Hello from /")
	//})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", roteador)

}

func heloHandler(w http.ResponseWriter, r *ninaRouter.NinaRequest) {
	fmt.Println("Request to /hello")
	jonas := r.PathValue("id")
	fmt.Fprint(w, "Hello from get")
	fmt.Printf("Hello from /hello/%s", jonas)
}
