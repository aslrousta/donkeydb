package paging

import (
	"errors"
	"io"
)

// New instantiates a new page-file.
func New(s io.ReadWriteSeeker, pageSize int) (*File, error) {
	if pageSize < 1 {
		return nil, errors.New("paging: invalid page size")
	}
	size, err := s.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}
	return &File{
		Stream:   s,
		PageSize: int64(pageSize),
		Pages:    size / int64(pageSize),
	}, nil
}

// Page is a sequence of bytes in the page-file.
type Page struct {
	Data  []byte
	Index int64
}

// File is a page-file over an underlying random-access stream.
type File struct {
	Stream   io.ReadWriteSeeker
	PageSize int64
	Pages    int64
}

// Read reads a page at the given index from the page-file.
func (f *File) Read(index int64) (*Page, error) {
	if index < 0 || index >= f.Pages {
		return nil, errors.New("paging: invalid page index")
	}
	offset := index * f.PageSize
	if _, err := f.Stream.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}
	page := &Page{
		Data:  make([]byte, f.PageSize),
		Index: index,
	}
	if _, err := f.Stream.Read(page.Data); err != nil {
		return nil, err
	}
	return page, nil
}

// Write writes a page back to the page-file.
func (f *File) Write(page *Page) error {
	offset := page.Index * f.PageSize
	if _, err := f.Stream.Seek(offset, io.SeekStart); err != nil {
		return err
	}
	if _, err := f.Stream.Write(page.Data); err != nil {
		return err
	}
	return nil
}

// Alloc allocates a new page at the end of the page-file.
func (f *File) Alloc() (*Page, error) {
	page := &Page{
		Data:  make([]byte, f.PageSize),
		Index: f.Pages,
	}
	if err := f.Write(page); err != nil {
		return nil, err
	}
	f.Pages++
	return page, nil
}
