package server

import (
	"fmt"
	"net"
)

type ConnectionHandler struct {
	id int
	conn net.Conn
	server *Server
}

func NewConnectionHandler(id int, conn net.Conn, server *Server) *ConnectionHandler {
	return &ConnectionHandler{id: id, conn: conn, server: server}
}

func (c *ConnectionHandler) ListenLoop() {
	tmp := make([]byte, 1024)
	for {
		_, err := c.conn.Read(tmp)
		//fmt.Println("Received message")
		if err != nil {
			fmt.Println(err)
			c.server.DeadClient(c.id)
			return
		}

		c.server.Retransmit(tmp, c.id)
	}
}

func (c *ConnectionHandler) Send(data []byte) {
	_, err := c.conn.Write(data)
	chk(err)
}

func (c *ConnectionHandler) GetId() int {
	return c.id
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
