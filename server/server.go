package server

import (
	"fmt"
	"net"
)

type Server struct {
	IP   string
	Port int
}

func CreateServer(ip string, port int) *Server {
	s := &Server{
		IP:   ip,
		Port: port,
	}
	return s
}

func (server *Server) StartServer() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.IP, server.Port))
	if err != nil {
		panic("network error")
		return
	}
	defer listener.Close()
	for {
		conn, e := listener.Accept()
		if e != nil {
			fmt.Println("network error accept", e)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	addr := conn.RemoteAddr()
	fmt.Println(addr)
	fmt.Println("got new request with")
}
