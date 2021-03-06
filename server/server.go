package server

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	IP          string
	Port        int
	OnlineUsers map[string]*User
	lock        sync.RWMutex
	MessageChan chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:          ip,
		Port:        port,
		OnlineUsers: make(map[string]*User),
		MessageChan: make(chan string),
	}
	return server
}

func (s *Server) ListenMessager() {
	for {
		msg := <-s.MessageChan

		s.lock.Lock()
		for _, user := range s.OnlineUsers {
			user.MessageChan <- msg
		}
		s.lock.Unlock()
	}
}

func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := user.Addr + "I am " + user.Name + msg
	s.MessageChan <- sendMsg
}

func (s *Server) connHandler(conn net.Conn) {
	//add user to online user list
	user := NewUser(conn, s)
	user.Online()

	//receive user's message then broadcast
	go func() {
		buff := make([]byte, 4096)
		for {
			n, err := conn.Read(buff)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil {
				fmt.Println("Conn Read err:", err)
				return
			}
			msg := string(buff[:n])
			user.SentMsgToAll(msg)
		}
	}()
	select {}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()

	go s.ListenMessager()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		go s.connHandler(conn)
	}

}
