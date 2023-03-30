package main

import (
	"log"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mandriota/tic-tac-toe/internal/client/ui"
	"github.com/mandriota/tic-tac-toe/pkg/world"
)

func main() {
	logFile, err := os.OpenFile("tic-tac-toe.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	l := log.New(logFile, "", log.LstdFlags)

	conn, err := net.Dial("tcp", ":88")
	if err != nil {
		l.Println(err)
	}
	defer conn.Close()

	conn.Write([]byte("test_room\n"))

	id := make([]byte, 1)
	conn.Read(id)

	wm := world.NewWorldMeta(world.NewWorldBase(15, 15, 5), conn, conn)

	if id[0] == 1 {
		wm.RemoteTryMove()
	}

	l.Println("starting UI...")

	p := tea.NewProgram(ui.NewModel(wm))
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}
