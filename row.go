package godb

import (
	"encoding/binary"
	"fmt"
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

	a := data[:unsafe.Sizeof(r.ID)]
	b := data[unsafe.Sizeof(r.ID) : unsafe.Sizeof(r.Username)+4]
	c := data[unsafe.Sizeof(r.Username)+4:]

	r.ID = binary.LittleEndian.Uint32(a)

	var username [32]byte
	copy(username[:], b)
	r.Username = username
	var email [256]byte
	copy(email[:], c)
	r.Email = email

	return nil
}
