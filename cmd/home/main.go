package main

import (
	"crypto/tls"
	"mime"

	"github.com/belak/home"
)

func main() {
	_ = mime.AddExtensionType(".gmi", "text/gemini")
	_ = mime.AddExtensionType(".gemini", "text/gemini")

	server := home.Server{
		TLS: &tls.Config{},
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
