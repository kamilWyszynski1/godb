package godb

import (
	"os"
)

const TableMaxRows = 1024

type Table struct {
	Pager       *Pager
	NumRows     uint32
	rootPageNum uint32
}

func OpenDB(filename string) (*Table, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return nil, err
	}

	fs, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return &Table{
		Pager: &Pager{
			f:          file,
			fileLength: uint32(fs.Size()),
			pages:      [TableMaxRows]*Row{},
		},
		NumRows: uint32(fs.Size()) / RowSize(),
	}, nil
}

func (t *Table) Close() error {
	return t.Pager.Close(t.NumRows)
}
