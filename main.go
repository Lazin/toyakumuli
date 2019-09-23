package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	iport = flag.Int("iport", 8181, "Ingestion port")
	qport = flag.Int("qport", 8282, "Qurey port")
)

func main() {

	fmt.Printf("Starting akumuli\nIngestion port: %d\nQuery port: %d\n", iport, qport)

	stop := make(chan struct{}, 1)
	out := make(chan Point, 1024)
	rsrv, err := NewRespServer(fmt.Sprintf(":%d", iport), out)
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

	tss := NewTSS()
	go func() {
		// Main ingestion loop, data from all connected clients goes here
		for point := range out {
			tss.Append(point)
			fmt.Println(point.series)
			fmt.Println(point.timestamp)
			fmt.Println(point.value)
		}
	}()

	<-stop
}
