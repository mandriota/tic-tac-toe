package world

import (
	"encoding/binary"
	"io"
	"sync"
)

type WorldMeta struct {
	*WorldBase

	sync.Mutex

	r io.Reader
	w io.Writer
}

func NewWorldMeta(wb *WorldBase, r io.Reader, w io.Writer) *WorldMeta {
	return &WorldMeta{
		WorldBase: wb,
		r:         r,
		w:         w,
	}
}

func (wm *WorldMeta) RemoteTryMove() bool {
	row := int32(0)
	col := int32(0)
	binary.Read(wm.r, binary.LittleEndian, &row)
	binary.Read(wm.r, binary.LittleEndian, &col)
	return wm.WorldBase.TryMove(row, col)
}

func (wm *WorldMeta) TryMove(row, col int32) bool {
	if !wm.Mutex.TryLock() {
		return false
	}

	binary.Write(wm.w, binary.LittleEndian, row)
	binary.Write(wm.w, binary.LittleEndian, col)
	wm.WorldBase.TryMove(row, col)

	go func() {
		wm.RemoteTryMove()
		wm.Mutex.Unlock()
	}()

	return true
}
