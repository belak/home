package main

import (
	"crypto/tls"
	"mime"
	"os"

	"github.com/belak/home"
)

func main() {
	_ = mime.AddExtensionType(".gmi", "text/gemini")
	_ = mime.AddExtensionType(".gemini", "text/gemini")

	server := home.Server{
		TLS:       &tls.Config{},
		Content:   os.DirFS("content"),
		Static:    os.DirFS("content/static"),
		Templates: os.DirFS("content/templates"),
	}

	cert, err := tls.LoadX509KeyPair("./server-cert.pem", "./server-key.pem")
	if err != nil {
		panic(err.Error())
	}
	server.TLS.Certificates = []tls.Certificate{cert}

	err = server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}
