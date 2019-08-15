package main

import (
	"./network"
)

func main() {
	server := network.NewServer()

	server.Start()
}