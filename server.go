package main

import (
	"github.com/orc/router"
	"net/http"
	"os"
)

func main() {
	println("Server started.")
	http.Handle("/", new(router.FastCGIServer))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./static/img"))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		println("Error listening: ", err.Error())
		os.Exit(1)
	}
}
