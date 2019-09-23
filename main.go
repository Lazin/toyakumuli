package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stop := make(chan struct{}, 1)
	out := make(chan Point, 1024)
	rsrv, err := NewRespServer("localhost:8182", out)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsrv.Close()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Println("signal received")
		stop <- struct{}{}
		close(out)
	}()

	go func() {
		// Main ingestion loop, data from all connected clients goes here
		for point := range out {
			fmt.Println(point.series)
			fmt.Println(point.timestamp)
			fmt.Println(point.value)
		}
	}()

	<-stop
}
