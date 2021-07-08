package godb

import (
	"encoding/binary"
	"fmt"
	"os"
	"unsafe"
)

type Row struct {
	ID       uint32
	Username [32]byte
	Email    [256]byte
}

func (r Row) String() string {
	return fmt.Sprintf("(%d %s %s)", r.ID, string(r.Username[:]), string(r.Email[:]))
}

func RowSize() uint32 {
	return 4 + 32 + 256
}

func (r *Row) Marshal() ([]byte, error) {
	b := make([]byte, 4)

	binary.LittleEndian.PutUint32(b, r.ID)
	for _, by := range r.Username {
		b = append(b, by)
	}
	for _, by := range r.Email {
		b = append(b, by)
	}
	return b, nil
}

func UnmarshalRow(data []byte, r *Row) error {
	fmt.Println(unsafe.Sizeof(r.ID))
	fmt.Println(unsafe.Sizeof(r.Username))
	a := data[:unsafe.Sizeof(r.ID)]
	b := data[unsafe.Sizeof(r.ID) : unsafe.Sizeof(r.Username)+4]
	c := data[unsafe.Sizeof(r.Username)+4:]

	fmt.Println(binary.LittleEndian.Uint32(a))
	fmt.Println(string(b))
	fmt.Println(string(c))

	r.ID = binary.LittleEndian.Uint32(a)

	var username [32]byte
	copy(username[:], b)
	r.Username = username
	var email [256]byte
	copy(email[:], c)
	r.Email = email

	return nil
}

const TableMaxRows = 1024

type Table struct {
	Pager   *Pager
	NumRows uint32
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
		NumRows: 0,
	}, nil
}

func (t *Table) Close() error {
	return t.Pager.Close(t.NumRows)
}
