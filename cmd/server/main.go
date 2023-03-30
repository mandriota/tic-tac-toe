package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/mandriota/tic-tac-toe/internal/server"
)

func main() {
	logFName := flag.String("l", "tic-tac-toe_server.log", "log file name")
	url := flag.String("url", ":4040", "The url to listen on")
	flag.Parse()

	logFile, err := os.OpenFile(*logFName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer logFile.Close()

	l := log.New(logFile, "", log.LstdFlags)

	listen, err := net.Listen("tcp", *url)
	if err != nil {
		l.Fatalln(err)
	}
	defer listen.Close()

	handler := server.NewHandler()

	l.Println("listenning...")

	for {
		conn, err := listen.Accept()
		if err != nil {
			l.Println(err)
		}
		go handler.Handle(conn)
		l.Println("connection accepted.")
	}
}
