package main

import (
	"github.com/kamuridesu/ip-syncer/internal/client"
	"github.com/kamuridesu/ip-syncer/internal/server"
	"os"
)

func main() {
	if os.Getenv("CLIENT") != "" {
		client.Start()
		return
	}
	handler, err := server.NewHandler("sqlite", "./data.db")
	if err != nil {
		panic(err)
	}
	server.Start(handler)
}
