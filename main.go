package main

import (
	"log"
	"net"
	"tcp-server/models"
)

func main() {

	s := models.NewServer()

	go s.Run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start the server: %s", err.Error())
	}
	defer listener.Close()
	log.Println("started tcp server on :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection %s", err.Error())
			continue
		}
		go s.NewClient(conn)
	}

}
