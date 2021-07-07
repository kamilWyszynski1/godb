package godb

import "fmt"

type ExecuteResult uint8

const (
	ExecuteSuccess ExecuteResult = iota + 1
	ExecuteTableFull
	ExecuteUnknown
)

func ExecuteStatement(s *Statement, t *Table) ExecuteResult {
	switch s.Type {
	case StatementInsert:
		if t.NumRows >= TableMaxRows {
			return ExecuteTableFull
		}
		t.Pager.pages = append(t.Pager.pages, s.RowToInsert)
		t.NumRows++

		return ExecuteSuccess
	case StatementSelect:
		for _, r := range t.Pager.pages {
			fmt.Print(r)
		}
	}
	return ExecuteUnknown
}
