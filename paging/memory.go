package paging

import (
	"errors"
	"io"
)

// Memory is an in-memory implementation of io.ReadWriteSeeker.
type Memory struct {
	Data   []byte
	Cursor int
}

// Read reads upto len(p) bytes into p from the stream.
func (m *Memory) Read(p []byte) (int, error) {
	if m.Cursor >= len(m.Data) {
		return 0, io.EOF
	}
	n := len(p)
	if n > 0 {
		if rem := len(m.Data) - m.Cursor; n > rem {
			n = rem
		}
		copy(p, m.Data[m.Cursor:m.Cursor+n])
	}
	return n, nil
}

// Write writes upto len(p) bytes from p to the stream.
func (m *Memory) Write(p []byte) (int, error) {
	n := len(p)
	if last := m.Cursor + n; last > len(m.Data) {
		data := make([]byte, last)
		copy(data, m.Data)
		m.Data = data
	}
	copy(m.Data[m.Cursor:], p)
	return n, nil
}

// Seek sets the offset for the next Read or Write.
func (m *Memory) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekCurrent:
		offset += int64(m.Cursor)
	case io.SeekEnd:
		offset = int64(len(m.Data)) - offset
	}
	if offset < 0 {
		return 0, errors.New("paging: invalid offset")
	}
	m.Cursor = int(offset)
	return offset, nil
}
