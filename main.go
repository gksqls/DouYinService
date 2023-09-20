package main

import (
	"DouYinService/route"
)

func main() {
	server := new(route.Server)
	server.Start()
}
