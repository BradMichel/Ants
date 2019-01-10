package main

import (
	"net/http"

	r "./routers"

	"github.com/gorilla/context"
)

func main() {

	//runtime.GOMAXPROCS(runtime.NumCPU())

	// http.HandleFunc("/", handler)
	router := r.InitRoutes()
	// http.Handle("/api/", router)

	http.Handle("/api/", context.ClearHandler(router))
	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("./files/data/"))))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static/build/"))))
	http.ListenAndServe("0.0.0.0:8080", nil)

}
