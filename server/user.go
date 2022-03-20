package server

import (
	"net"
)

type User struct {
	Name        string
	Addr        string
	MessageChan chan string
	conn        net.Conn
	server      *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:        userAddr,
		Addr:        userAddr,
		MessageChan: make(chan string),
		conn:        conn,
		server:      server,
	}

	go user.ListenMessage()
	return user
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.MessageChan
		u.conn.Write([]byte(msg))
	}
}

func (u *User) Online() {
	u.server.lock.Lock()
	u.server.OnlineUsers[u.Name] = u
	u.server.lock.Unlock()

	//broadcast user online
	u.server.BroadCast(u, "online")
}

func (u *User) Offline() {
	u.server.lock.Lock()
	delete(u.server.OnlineUsers, u.Name)
	u.server.lock.Unlock()

	//broadcast user online
	u.server.BroadCast(u, "offline")
}

func (u *User) SentMsgToAll(msg string) {
	u.server.BroadCast(u, msg)
}

func (u *User) sentMsgToSelf(msg string) {
	u.conn.Write([]byte(msg))
}
