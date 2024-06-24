package main

import (
	"github.com/belak/home"
)

func main() {
	server := home.Server{}

	err := server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}
