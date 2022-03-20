package server

import (
	"fmt"
	"net"
)

type Server struct{
	IP string
	Port int
}

func NewServer(ip string, port int) *Server{
	server:=&Server{
		IP: ip,
		Port:port,
	}
	return server
}

func (s *Server) connHandler(conn net.Conn){
	fmt.Println("I am a handler")
}

func(s *Server) Start(){
	listener,err:=net.Listen("tcp",fmt.Sprintf("%s:%d", s.IP,s.Port))
	if err!=nil{
		fmt.Println("net.Listen err:",err)
		return
	}
	defer listener.Close()

	for{
		conn,err:=listener.Accept()
		if err!=nil{
			fmt.Println("listener accept err:",err)
			continue
		}
		go s.connHandler(conn)
	}

}

