package main

import (
	"github.com/wcytk/trainThroughBlockchain/network"
)

func main() {
	server := network.NewServer()

	server.Start()
}