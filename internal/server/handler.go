package server

import (
	"bufio"
	"io"
	"log"
	"net"
)

type Room struct {
	x chan net.Conn
	o chan net.Conn
}

type Handler struct {
	rooms map[string]*Room
}

func NewHandler() *Handler {
	return &Handler{rooms: make(map[string]*Room, 1024)}
}

func (h *Handler) Handle(conn net.Conn) {
	defer conn.Close()

	sc := bufio.NewScanner(conn)

	for sc.Scan() {
		name := sc.Text()
		log.Println(name)
		room, ok := h.rooms[name]
		if !ok {
			room = &Room{
				x: make(chan net.Conn, 2),
				o: make(chan net.Conn, 2),
			}
			h.rooms[name] = room
		}

		other := net.Conn(nil)

		if !ok {
			room.x <- conn
			other = <-room.o

			conn.Write([]byte{0})
		} else {
			conn.Write([]byte{1})
			room.o <- conn
			other = <-room.x
		}

		io.Copy(conn, other)

		delete(h.rooms, name)
	}
}
