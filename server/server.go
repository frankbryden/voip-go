package server

import (
	"fmt"
	"net"
)

type Server struct {
	ln net.Listener
	connections []*ConnectionHandler
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Listen() {
	ln, err := net.Listen("tcp", ":8091")
	fmt.Println("Listening on 8080")
	if err != nil {
		fmt.Println(err)
	}
	id := 0
	for {
		conn, err := ln.Accept()
		fmt.Println("Connection received!")
		if err != nil {
			// handle error
			fmt.Println(err)
		}
		//Create a connection handler - passing in a unique id and a handle back to this server
		//the handle is needed so data can be retransmitted to all active connections
		handler := NewConnectionHandler(id, conn, s)
		s.connections = append(s.connections, handler)
		go handler.ListenLoop()
		id++
	}
}

func (s *Server) Retransmit(data []byte, id int) {
	//fmt.Printf("Retransmit from client %d\n", id)
	for _, connHandler := range s.connections {
		if connHandler.GetId() == id {
			continue
		}
		connHandler.Send(data)
	}
}

func (s *Server) DeadClient(id int) {
	//Remove a dead client from our list of active connections
	s.connections = append(s.connections[:id], s.connections[id+1:]...)
}