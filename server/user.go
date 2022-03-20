package server

import (
	"net"
)

type User struct {
	Name        string
	Addr        string
	MessageChan chan string
	conn        net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:        userAddr,
		Addr:        userAddr,
		MessageChan: make(chan string),
		conn:        conn,
	}

	go user.ListenMessage()
	return  user
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.MessageChan
		u.conn.Write([]byte(msg))
	}
}
