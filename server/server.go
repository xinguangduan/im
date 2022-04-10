package server

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	IP        string
	Port      int
	OnlineMap map[string]*User
	MapLock   sync.RWMutex
	Message   chan string // server channel
}

func CreateServer(ip string, port int) *Server {
	s := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return s
}

func (s *Server) StartServer() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		panic("network error")
		return
	}
	defer listener.Close()
	// start listen
	go s.ListenMessage()

	for {
		conn, e := listener.Accept()
		if e != nil {
			fmt.Println("network error accept", e)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.MapLock.Lock()
		for key, user := range s.OnlineMap {
			fmt.Println(key, user.Name)
			user.Channel <- msg
		}
		s.MapLock.Unlock()
	}
}

func (s *Server) broadCastMessage(u *User, msg string) {
	bMsg := "[" + u.Addr + "]" + msg
	s.Message <- bMsg
}

func (s *Server) handle(conn net.Conn) {
	addr := conn.RemoteAddr()
	fmt.Println(addr)
	fmt.Println("got new request with")
	u := CreateUser(conn)
	s.MapLock.Lock()
	s.OnlineMap[u.Name] = u
	s.MapLock.Unlock()
	// send message to User
	s.broadCastMessage(u, "already online")
	// block current process
	select {}
}
