package paging

import (
	"errors"
	"io"
)

// New instantiates a new page-file.
func New(rws io.ReadWriteSeeker, pageSize int) (*File, error) {
	if pageSize < 1 {
		return nil, errors.New("paging: invalid page size")
	}
	size, err := rws.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}
	return &File{
		rws:      rws,
		pageSize: int64(pageSize),
		pages:    size / int64(pageSize),
	}, nil
}

// Page is a sequence of bytes in the page-file.
type Page struct {
	Data  []byte
	Index int64
}

// File is a page-file over an underlying random-access stream.
type File struct {
	rws      io.ReadWriteSeeker
	pageSize int64
	pages    int64
}

// Read reads a page at the given index from the page-file.
func (f *File) Read(index int64) (*Page, error) {
	if index < 0 || index >= f.pages {
		return nil, errors.New("paging: invalid page index")
	}
	offset := index * f.pageSize
	if _, err := f.rws.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}
	page := &Page{
		Data:  make([]byte, f.pageSize),
		Index: index,
	}
	if _, err := f.rws.Read(page.Data); err != nil {
		return nil, err
	}
	return page, nil
}

// Write writes a page back to the page-file.
func (f *File) Write(page *Page) error {
	offset := page.Index * f.pageSize
	if _, err := f.rws.Seek(offset, io.SeekStart); err != nil {
		return err
	}
	if _, err := f.rws.Write(page.Data); err != nil {
		return err
	}
	return nil
}

// Alloc allocates a new page at the end of the page-file.
func (f *File) Alloc() (*Page, error) {
	page := &Page{
		Data:  make([]byte, f.pageSize),
		Index: f.pages,
	}
	if err := f.Write(page); err != nil {
		return nil, err
	}
	f.pages++
	return page, nil
}
