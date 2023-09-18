package main

import (
	"DouYinService/route"
)

func main() {
	server := &route.Server{
		Bind: ":9000", DouYinUrl: "https://live.douyin.com/26242350553",
	}
	server.Start()
}
