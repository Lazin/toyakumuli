package main

import "fmt"

func main() {
	rsrv, err := NewRespServer("localhost:8181")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsrv.Close()
}
