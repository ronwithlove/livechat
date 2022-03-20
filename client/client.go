package main

import (
	"fmt"
	"net"
	"time"
)

type Client struct {
	ServerIP   string
	ServerPort int
	conn       net.Conn
}

func NewClient(serverIP string, serverPort int) *Client {
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net.Dail err: ", err)
		return nil
	}

	client.conn = conn
	return client
}

func (c *Client) SendToAllUsers(msg string) {
	if len(msg) != 0 {
		_, err := c.conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("conn write err:", err)
		}
	}
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("connect to server failed.")
		return
	}
	fmt.Println("Successfully connected to the server")

	//test send msg to server
	go func() {
		for {
			time.Sleep(5 * time.Second)
			client.SendToAllUsers("I am from " + client.conn.RemoteAddr().String())
		}
	}()

	//get msg from server
	var buff [512]byte
	for {
		_, err := client.conn.Read(buff[0:])
		if err != nil {
			return
		}
		fmt.Println(string(buff[0:]))
	}
}
