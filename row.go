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
	b := make([]byte, RowSize())

	binary.LittleEndian.PutUint32(b, r.ID)
	for i, by := range r.Username {
		b[4+i] = by
	}
	for i, by := range r.Email {
		b[36+i] = by
	}
	return b, nil
}

func UnmarshalRow(data []byte, r *Row) error {
	idSize := unsafe.Sizeof(r.ID)
	usernameSize := unsafe.Sizeof(r.Username)
	a := data[:idSize]
	b := data[idSize : idSize+usernameSize]
	c := data[idSize+usernameSize:]

	r.ID = binary.LittleEndian.Uint32(a)

	var username [32]byte
	copy(username[:], b)
	r.Username = username
	var email [256]byte
	copy(email[:], c)
	r.Email = email

	return nil
}
