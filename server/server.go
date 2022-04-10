package server

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	IP        string
	Port      int
	OnlineMap map[string]*User

	Message chan string // server channel
}

var mapLock sync.RWMutex

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
		mapLock.Lock()
		for _, client := range s.OnlineMap {
			client.Channel <- msg
		}
		mapLock.Unlock()
	}
}

func (s *Server) BroadCastMessage(u *User, msg string) {
	bMsg := "[" + u.Addr + "]" + msg
	s.Message <- bMsg
}

func (s *Server) handle(conn net.Conn) {
	addr := conn.RemoteAddr()
	fmt.Println("new connection", addr.String())
	u := CreateUser(conn, *s)
	//s.MapLock.Lock()
	//s.OnlineMap[u.Name] = u
	//s.MapLock.Unlock()
	// send message to User
	//s.BroadCastMessage(u, "already online")
	u.Online()
	// read client data and process
	go s.handleClientMessage(conn, u)
	// block current process
	select {}
}

func (s *Server) handleClientMessage(conn net.Conn, u *User) {
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("read data from client error", err)
			break
		}
		if n == 0 {
			//s.BroadCastMessage(u, "下线啦")
			u.Offline()
			return
		}
		uMsg := string(buf)
		//s.BroadCastMessage(u, uMsg)
		u.HandleMessage(uMsg)
	}
}

func (s *Server) PushUserToMap(u *User) {
	mapLock.Lock()
	s.OnlineMap[u.Name] = u
	mapLock.Unlock()
}
func (s *Server) RemoveUserFromMap(u *User) {
	mapLock.Lock()
	delete(s.OnlineMap, u.Name)
	mapLock.Unlock()
}
