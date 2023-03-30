package world

type World interface {
	Rows() int32
	Cols() int32
	Look(row, col int32) byte
	TryMove(row, col int32) bool
	Winner() byte
}

const (
	dirU = iota
	dirL
	dirD
	dirR
	dirUL
	dirUR
	dirDL
	dirDR
)

type ceil struct {
	dir [8]*ceil

	mark byte
}

type WorldBase struct {
	board [][]ceil
	goal  int32 // winning line length
	last  *ceil // last move
	mark  byte  // current mark
}

func NewWorldBase(row, col, goal int32) *WorldBase {
	w := &WorldBase{
		board: make([][]ceil, row),
		goal:  goal,
		mark:  'x',
	}

	for i := range w.board {
		w.board[i] = make([]ceil, col)
	}

	for i := range w.board {
		for j := range w.board[i] {
			w.board[i][j].mark = 0
			if i > 0 {
				w.board[i][j].dir[dirU] = &w.board[i-1][j]
			}
			if j > 0 {
				w.board[i][j].dir[dirL] = &w.board[i][j-1]
			}
			if i < len(w.board)-1 {
				w.board[i][j].dir[dirD] = &w.board[i+1][j]
			}
			if j < len(w.board[i])-1 {
				w.board[i][j].dir[dirR] = &w.board[i][j+1]
			}
			if i > 0 && j > 0 {
				w.board[i][j].dir[dirUL] = &w.board[i-1][j-1]
			}
			if i > 0 && j < len(w.board[i])-1 {
				w.board[i][j].dir[dirUR] = &w.board[i-1][j+1]
			}
			if i < len(w.board)-1 && j > 0 {
				w.board[i][j].dir[dirDL] = &w.board[i+1][j-1]
			}
			if i < len(w.board)-1 && j < len(w.board[i])-1 {
				w.board[i][j].dir[dirDR] = &w.board[i+1][j+1]
			}
		}
	}

	return w
}

func (w *WorldBase) Rows() int32 {
	return int32(len(w.board))
}

func (w *WorldBase) Cols() int32 {
	if w.Rows() != 0 {
		return int32(len(w.board[0]))
	}
	return 0
}

func (w *WorldBase) Look(row, col int32) byte {
	return w.board[row][col].mark
}

func (w *WorldBase) TryMove(row, col int32) bool {
	if w.board[row][col].mark != 0 {
		return false
	}

	w.board[row][col].mark = w.mark
	w.last = &w.board[row][col]
	if w.mark == 'x' {
		w.mark = 'o'
	} else {
		w.mark = 'x'
	}

	return true
}

func (w WorldBase) checkLine(dir int) (cnt int32) {
	for c := w.last; c.dir[dir] != nil && c.dir[dir].mark == c.mark; c = c.dir[dir] {
		cnt++
	}

	return
}

func (w WorldBase) Winner() byte {
	if w.last == nil {
		return 0
	}

	udC := 1 + w.checkLine(dirU) + w.checkLine(dirD)
	if udC >= w.goal {
		return w.last.mark
	}

	lrC := 1 + w.checkLine(dirL) + w.checkLine(dirR)
	if lrC >= w.goal {
		return w.last.mark
	}

	uldrC := 1 + w.checkLine(dirUL) + w.checkLine(dirDR)
	if uldrC >= w.goal {
		return w.last.mark
	}

	urdlC := 1 + w.checkLine(dirUR) + w.checkLine(dirDL)
	if urdlC >= w.goal {
		return w.last.mark
	}

	return 0
}
