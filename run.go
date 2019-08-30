package main

import (
	"github.com/wcytk/BlockTrain/network"
)

func main() {
	server := network.NewServer()

	server.Start()
}