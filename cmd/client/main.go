package main

import (
	"flag"
	"log"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mandriota/tic-tac-toe/internal/client/ui"
	"github.com/mandriota/tic-tac-toe/pkg/world"
)

func main() {
	logFName := flag.String("l", "tic-tac-toe_client.log", "log file name")
	url := flag.String("u", ":4040", "server url")
	room := flag.String("r", "default", "room name")
	rows := flag.Int("x", 15, "number of rows")
	cols := flag.Int("y", 15, "number of columns")
	goal := flag.Int("g", 5, "goal")
	flag.Parse()

	logFile, err := os.OpenFile(*logFName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer logFile.Close()

	l := log.New(logFile, "", log.LstdFlags)

	conn, err := net.Dial("tcp", *url)
	if err != nil {
		l.Fatalln(err)
	}
	defer conn.Close()

	conn.Write([]byte(*room + "\n"))

	id := make([]byte, 1)
	conn.Read(id)

	wm := world.NewWorldMeta(world.NewWorldBase(int32(*rows), int32(*cols), int32(*goal)), conn, conn)

	if id[0] == 1 {
		wm.RemoteTryMove()
	}

	l.Println("starting UI...")

	p := tea.NewProgram(ui.NewModel(wm))
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}
