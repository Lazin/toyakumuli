package main

import (
	"bufio"
	"fmt"
	"net"
)

// RespServer is a RESP tcp server connection
type RespServer struct {
	listener net.Listener
	addr     string
	out      chan []byte
	done     chan struct{}
}

// NewRespServer creates new RESP server
func NewRespServer(addr string) (*RespServer, error) {
	var r RespServer
	l, err := net.Listen("tcp4", addr)
	if err != nil {
		return nil, err
	}
	r.listener = l
	r.addr = addr
	r.out = make(chan []byte, 1024)
	r.done = make(chan struct{})
	go func() {
		for {
			conn, err := r.listener.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}
			go r.processInput(conn)
		}
	}()
	return &r, err
}

func (r *RespServer) processInput(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		select {
		case r.out <- scanner.Bytes():
		case <-r.done:
			break
		}
	}
}

// Close shots down RESP server
func (r *RespServer) Close() {
	r.listener.Close()
	r.done <- struct{}{}
}
