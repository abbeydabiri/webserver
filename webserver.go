package main

import (
	"flag"
	"net/http"
)

func main() {

	var cPort string
	flag.StringVar(&cPort, "port", "8080", "Default Port")
	flag.Parse()

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(""))))

	println("Strting Server @ port: " + cPort)
	if err := http.ListenAndServe(":"+cPort, nil); err != nil {
		println("Error Starting Webserver: " + err.Error())
	}
}
