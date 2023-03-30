package ui

import (
	"bytes"

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

type Model struct {
	world *world.WorldMeta
	win   byte
	row   int32
	col   int32
	buf   *bytes.Buffer

	altScreen bool
}

func NewModel(w *world.WorldMeta) Model {
	return Model{
		world: w,
		row:   (w.Rows() - 1) / 2,
		col:   (w.Cols() - 1) / 2,
		buf:   bytes.NewBuffer(nil),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.row != m.world.Rows()-1 {
				m.row++
			}
		case "right":
			if m.col != m.world.Cols()-1 {
				m.col++
			}
		case "left":
			if m.col != 0 {
				m.col--
			}
		case "enter":
			if !m.world.TryMove(m.row, m.col) {
				return m, nil
			}
			m.win = m.world.Winner()

			return m, func() tea.Msg {
				m.world.Mutex.Lock()
				m.win = m.world.Winner()
				m.world.Mutex.Unlock()
				return true
			}
		case " ":
			m.altScreen = !m.altScreen
			if m.altScreen {
				return m, tea.EnterAltScreen
			}

			return m, tea.ExitAltScreen
		}
	case bool:
	}

	return m, nil
}

func (m Model) View() string {
	m.buf.Reset()

	for i := int32(0); i < m.world.Rows(); i++ {
		if i != 0 {
			m.buf.WriteByte('\n')
		}
		m.buf.WriteByte(' ')

		for j := int32(0); j < m.world.Cols(); j++ {
			s := gridStyle.Render("·")

			if b := m.world.Look(i, j); b != 0 {
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
