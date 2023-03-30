package main

import (
	"bytes"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mandriota/tic-tac-toe/pkg/world"
)

const (
	clrFuxia = lipgloss.Color("#FF006F")
)

var (
	selectionStyle = lipgloss.NewStyle().
			Foreground(clrFuxia).
			Bold(true)
	boardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(clrFuxia)
	gridStyle = lipgloss.NewStyle().Faint(true)
)

func main() {
	p := tea.NewProgram(newModel(world.NewWorld(15, 15, 5)))
	if _, err := p.Run(); err != nil {
		log.Println(err)
	}
}

type model struct {
	w   *world.World
	win byte
	row int
	col int
	buf *bytes.Buffer

	altScreen bool
}

func newModel(w *world.World) model {
	return model{
		w:   w,
		row: (w.Rows() - 1) / 2,
		col: (w.Cols() - 1) / 2,
		buf: bytes.NewBuffer(nil),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.win != 0 {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch s := msg.String(); s {
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.row != 0 {
				m.row--
			}
		case "down":
			if m.row != m.w.Rows()-1 {
				m.row++
			}
		case "right":
			if m.col != m.w.Cols()-1 {
				m.col++
			}
		case "left":
			if m.col != 0 {
				m.col--
			}
		case "enter":
			m.w.Move(m.row, m.col)
			m.win = m.w.Winner()
		case " ":
			m.altScreen = !m.altScreen
			if m.altScreen {
				return m, tea.EnterAltScreen
			}

			return m, tea.ExitAltScreen
		}
	}

	return m, nil
}

func (m model) View() string {
	m.buf.Reset()

	for i := 0; i < m.w.Rows(); i++ {
		if i != 0 {
			m.buf.WriteByte('\n')
		}
		m.buf.WriteByte(' ')

		for j := 0; j < m.w.Cols(); j++ {
			s := gridStyle.Render("·")

			if b := m.w.Look(i, j); b != 0 {
				s = string(b)
			} else if i == m.row || j == m.col {
				s = "•"
			}

			if i == m.row && j == m.col {
				s = selectionStyle.Render(s)
			}

			m.buf.WriteString(s)
			m.buf.WriteByte(' ')
		}
	}

	switch m.win {
	case 'x':
		return boardStyle.Render(m.buf.String()) + "\n X wins!"
	case 'o':
		return boardStyle.Render(m.buf.String()) + "\n O wins!"
	}

	return boardStyle.Render(m.buf.String())
}
