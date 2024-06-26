package main

import (
	"github.com/belak/home"
	"github.com/belak/home/internal"
)

func main() {
	logger, err := internal.NewLogger()
	if err != nil {
		panic(err.Error())
	}

	server := home.NewServer(home.ServerConfig{
		Logger:   logger,
		BindAddr: ":8080",
	})

	err = server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}
