package main

import (
	"os"

	"github.com/kamuridesu/ip-syncer/internal/client"
	"github.com/kamuridesu/ip-syncer/internal/hosts"
	"github.com/kamuridesu/ip-syncer/internal/server"
)

func main() {
	if os.Getenv("CLIENT") != "" {
		client.Start()
		return
	}
	const HostsPath = "/etc/hosts"
	hostFile, err := hosts.ReadHostsFile(HostsPath)
	if err != nil {
		panic(err)
	}
	handler, err := server.NewHandler(hostFile)
	if err != nil {
		panic(err)
	}
	server.Start(handler)
}
