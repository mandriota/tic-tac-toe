package main

import (
	"log"
	"net"

	"github.com/mandriota/tic-tac-toe/internal/server"
)

func main() {
	listen, err := net.Listen("tcp", ":88")
	if err != nil {
		log.Fatalln(err)
	}
	defer listen.Close()

	handler := server.NewHandler()

	log.Println("listenning...")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
		}
		go handler.Handle(conn)
		log.Println("connection accepted.")
	}
}
