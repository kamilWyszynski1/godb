package godb

import (
	"fmt"
	"strings"
)

type Status int8

const (
	UnrecognizedCommand Status = iota + 1
	SuccessfulExit
)

func Parse(input string) Status {
	switch input {
	case ".exit":
		return SuccessfulExit
	default:
		return UnrecognizedCommand
	}
}

type (
	StatementType int8
	PrepareStatus int8
)

const (
	StatementInsert StatementType = iota + 1
	StatementSelect

	PrepareSuccess PrepareStatus = iota + 1
	PrepareUnrecognizedStatement
	PrepareSyntaxError
)

type Statement struct {
	// Type defines what type of operation we will be doing
	Type StatementType
	// Status describes how preparing went
	Status      PrepareStatus
	RowToInsert *Row
}

func PrepareStatement(input string) *Statement {
	if strings.HasPrefix(input, "insert") {
		s := &Statement{Status: PrepareSuccess, Type: StatementInsert, RowToInsert: &Row{}}
		var username, email string
		if argsAssinged, err := fmt.Sscanf(input, "insert %d %v %v", &s.RowToInsert.ID, &username, &email); err != nil {
			return &Statement{Status: PrepareSyntaxError}
		} else if argsAssinged > 3 {
			return &Statement{Status: PrepareSyntaxError}
		}
		copy(s.RowToInsert.Username[:], username)
		copy(s.RowToInsert.Email[:], email)

		return s
	}

	if input == "select" {
		return &Statement{Status: PrepareSuccess, Type: StatementSelect}
	}
	return &Statement{Status: PrepareUnrecognizedStatement}

}
