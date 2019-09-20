package main

import (
	"net"
)

// RespServer is a RESP tcp server connection
type RespServer struct {
	listener net.Listener
	addr     string
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
	return &r, err
}

// Close shots down RESP server
func (r *RespServer) Close() {
	r.listener.Close()
}
