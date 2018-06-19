package WEBSERVER

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"os"
	"strings"

	"golang.org/x/mobile/asset"
)

var flagPORT, flagOS string

func main() {
	flag.StringVar(&flagOS, "os", "", "Default os")
	flag.StringVar(&flagPORT, "port", "8080", "Default Port")
	flag.Parse()
	startWebserver()
}

func START(os, port string) {
	flagOS = os
	flagPORT = port
	go startWebserver()
}

func startWebserver() {
	http.HandleFunc("/", siteHandler)
	println("Starting Server @ port: " + flagPORT)
	if err := http.ListenAndServe(":"+flagPORT, nil); err != nil {
		println("Error Starting Webserver: " + err.Error())
	}
}

func siteHandler(httpRes http.ResponseWriter, httpReq *http.Request) {
	urlPath := strings.Split(httpReq.URL.String()[1:], "?")[0]
	urlPath = strings.Replace(urlPath, "//", "/", -1)

	switch flagOS {
	case "ios", "android":
		if _, err := getAsset(urlPath); err != nil {
			urlPath = "index.html"
		}
	default:
		if _, err := os.Stat(urlPath); os.IsNotExist(err) {
			urlPath = "index.html"
		}
	}

	httpRes.Header().Set("Cache-Control", "max-age=3600, must-revalidate")
	httpRes.Header().Set("Pragma", "no-cache")
	httpRes.Header().Set("Expires", "0")

	if dataBytes, err := getAsset(urlPath); err == nil {
		httpRes.Header().Add("Content-Type", getContentType(urlPath))
		if !strings.Contains(httpReq.Header.Get("Accept-Encoding"), "gzip") {
			httpRes.Write(dataBytes)
			return
		}
		gzipWrite(dataBytes, httpRes)
	}
}

func getAsset(filename string) (assetByte []byte, assetError error) {
	assetByte = nil
	assetError = nil

	if strings.HasSuffix(filename, "/") {
		assetError = fmt.Errorf("directory listing forbidden")
	} else {
		switch flagOS {
		case "ios", "android":
			if f, errOpen := asset.Open(filename); errOpen == nil {
				defer f.Close()
				assetByte, assetError = ioutil.ReadAll(f)
			}
		default:
			assetByte, assetError = ioutil.ReadFile(filename)
		}
	}
	return
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipWrite(dataBytes []byte, httpRes http.ResponseWriter) {
	httpRes.Header().Set("Content-Encoding", "gzip")
	gzipHandler := gzip.NewWriter(httpRes)
	defer gzipHandler.Close()
	httpResGzip := gzipResponseWriter{Writer: gzipHandler, ResponseWriter: httpRes}
	httpResGzip.Write(dataBytes)
}

func getContentType(filename string) (contentType string) {
	// contentType = "text/plain; charset=utf-8"
	contentType = "text/html"
	switch {
	case strings.HasSuffix(filename, ".apk"):
		contentType = "application/vnd.android.package-archive"

	case strings.HasSuffix(filename, ".js"):
		contentType = "application/javascript"
	case strings.HasSuffix(filename, ".json"):
		contentType = "application/json"
	case strings.HasSuffix(filename, ".pdf"):
		contentType = "application/pdf"
	case strings.HasSuffix(filename, ".zip"):
		contentType = "application/zip"

	case strings.HasSuffix(filename, ".xls"):
		contentType = "application/vnd.ms-excel"
	case strings.HasSuffix(filename, ".xlsx"):
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

	case strings.HasSuffix(filename, ".html"):
		contentType = "text/html"
	case strings.HasSuffix(filename, ".css"):
		contentType = "text/css"

	case strings.HasSuffix(filename, ".doc"):
		contentType = "application/msword"
	case strings.HasSuffix(filename, ".docx"):
		contentType = "application/msword"

	case strings.HasSuffix(filename, ".png"):
		contentType = "image/png"
	case strings.HasSuffix(filename, ".jpg"),
		strings.HasSuffix(filename, ".jpeg"):
		contentType = "image/jpeg"
	case strings.HasSuffix(filename, ".gif"):
		contentType = "image/gif"
	case strings.HasSuffix(filename, ".svg"):
		contentType = "image/svg+xml"

	case strings.HasSuffix(filename, ".mp4"):
		contentType = "video/mp4"
	case strings.HasSuffix(filename, ".webm"):
		contentType = "video/webm"
	case strings.HasSuffix(filename, ".ogg"):
		contentType = "video/ogg"
	case strings.HasSuffix(filename, ".mp3"):
		contentType = "audio/mp3"
	case strings.HasSuffix(filename, ".wav"):
		contentType = "audio/wav"
	}
	return
}
