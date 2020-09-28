package main

import (
	"crud-using-chi/app"
	"fmt"
	"log"
	"net/http"
)

func main() {
	r, conf := app.Initialize()
	port := conf.GetString("app.port")
	fmt.Printf("\nServer started. Listening to port %v.\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, r))
}
