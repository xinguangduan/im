package server

import (
	"fmt"
	"net"
)

type User struct {
	Name    string
	Addr    string
	Channel chan string
	Conn    net.Conn
}

func CreateUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	u := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		Conn:    conn,
	}
	go u.ListenMessage()

	return u
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.Channel
		fmt.Println(msg)
		// write message to client
		u.Conn.Write([]byte(msg + "\n"))
	}
}
