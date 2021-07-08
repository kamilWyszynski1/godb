package godb

import (
	"fmt"
	"math/rand"
	"testing"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func prepRow(id uint32) *Row {
	var username [32]byte
	copy(username[:], RandStringRunes(32))

	var email [256]byte
	copy(email[:], RandStringRunes(256))
	return &Row{
		ID:       id,
		Username: username,
		Email:    email,
	}
}

func TestMarshalRow(t *testing.T) {
	r := prepRow(4294967295)

	data, err := r.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var r1 Row
	if err := UnmarshalRow(data, &r1); err != nil {
		t.Fatal(err)
	}
}

func TestWholeFlow(t *testing.T) {
	table, err := OpenDB("test.godb")
	if err != nil {
		t.Fatal(err)
	}
	for i := uint32(0); i < 10; i++ {
		r := prepRow(i)
		fmt.Printf("%d:%s", i, r)
		if err := ExecuteStatement(&Statement{
			Type:        StatementInsert,
			Status:      0,
			RowToInsert: r,
		}, table); err != nil {
			t.Fatal(err)
		}
	}

	if err := table.Close(); err != nil {
		t.Fatal(err)
	}

	if err := ExecuteStatement(&Statement{Type: StatementSelect}, table); err != nil {
		t.Fatal(err)
	}
}
