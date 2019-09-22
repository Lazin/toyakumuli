package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

const layout = "+20060102T150405.999999999"

// RespServer is a RESP tcp server connection
type RespServer struct {
	listener net.Listener
	addr     string
	out      chan<- Point
	done     chan struct{}
}

// NewRespServer creates new RESP server
func NewRespServer(addr string, out chan<- Point) (*RespServer, error) {
	var r RespServer
	l, err := net.Listen("tcp4", addr)
	if err != nil {
		return nil, err
	}
	r.listener = l
	r.addr = addr
	r.out = out
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
	var ix int = 0
	var point Point
	for scanner.Scan() {
		if ix%3 == 0 {
			sname := scanner.Text()
			point.series = sname[1:len(sname)]
		} else if ix%3 == 1 {
			time, err := time.Parse(layout, scanner.Text())
			if err != nil {
				// Ignoring errors and possibility of incomplete write
				conn.Write([]byte("!RESP error"))
				conn.Close()
			}
			point.timestamp = time
		} else if ix%3 == 2 {
			// TODO: Parse value
			point.value = 3.14159
			r.out <- point
		}
		select {
		case <-r.done:
			break
		}
		ix++
	}
}

// Close shots down RESP server
func (r *RespServer) Close() {
	r.listener.Close()
	r.done <- struct{}{}
}
