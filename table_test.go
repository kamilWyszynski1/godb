package godb

import (
	"testing"
)

func TestMarshalRow(t *testing.T) {
	var username [32]byte
	copy(username[:], "test uername")

	var email [256]byte
	copy(email[:], "test email")
	r := Row{
		ID:       4294967295,
		Username: username,
		Email:    email,
	}

	data, err := MarshalRow(r)
	if err != nil {
		t.Fatal(err)
	}
	var r1 Row
	if err := UnmarshalRow(data, &r1); err != nil {
		t.Fatal(err)
	}
}
