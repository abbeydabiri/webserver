package main

import (
	"flag"
	"net/http"

	"os"
	"strings"
)

func main() {
	var cPort string
	flag.StringVar(&cPort, "port", "8080", "Default Port")

	http.HandleFunc("/", siteHandler)
	println("Starting Server @ port: " + cPort)
	if err := http.ListenAndServe(":"+cPort, nil); err != nil {
		println("Error Starting Webserver: " + err.Error())
	}
}

func siteHandler(httpRes http.ResponseWriter, httpReq *http.Request) {
	urlPath := strings.Split(httpReq.URL.String()[1:], "?")[0]
	urlPath = strings.Replace(urlPath, "//", "/", -1)
	if _, err := os.Stat(urlPath); os.IsNotExist(err) {
		urlPath = "index.html"
	}
	http.ServeFile(httpRes, httpReq, urlPath)
}
