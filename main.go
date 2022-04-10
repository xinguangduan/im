package main

import (
	"flag"
	"fmt"
	"github.com/xinguangduan/im/server"
)

func main() {
	var ip string
	var port int
	flag.StringVar(&ip, "h", "127.0.0.1", "")
	flag.IntVar(&port, "p", 19800, "")
	flag.Parse()
	fmt.Println("start server...")
	s := server.CreateServer(ip, port)
	s.StartServer()
}
