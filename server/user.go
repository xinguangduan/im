package server

import (
	"net"
)

type User struct {
	Name    string
	Addr    string
	Channel chan string
	Conn    net.Conn
	Server  Server
}

func CreateUser(conn net.Conn, s Server) *User {
	userAddr := conn.RemoteAddr().String()
	u := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		Conn:    conn,
		Server:  s,
	}
	go u.ListenMessage()
	return u
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.Channel
		// write message to client
		u.Conn.Write([]byte(msg + "\n"))
	}
}

func (u *User) Online() {
	u.Server.PushUserToMap(u)
	u.Server.BroadCastMessage(u, u.Name+"已上线")
}

func (u *User) Offline() {
	u.Server.RemoveUserFromMap(u)
	u.Server.BroadCastMessage(u, u.Name+"已下线")
}

func (u *User) HandleMessage(msg string) {
	u.Server.BroadCastMessage(u, msg)
}
