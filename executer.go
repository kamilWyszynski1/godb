package godb

import (
	"errors"
	"fmt"
)

var (
	ErrExecuteTableFull = errors.New("table is full")
	ErrExecuteUnknown   = errors.New("unknown statement")
)

func ExecuteStatement(s *Statement, t *Table) error {
	switch s.Type {
	case StatementInsert:
		if t.NumRows >= TableMaxRows {
			return ErrExecuteTableFull
		}
		t.Pager.pages[t.NumRows] = s.RowToInsert
		t.NumRows++

		return nil
	case StatementSelect:
		for i := uint32(0); i < t.NumRows; i++ {
			row, err := t.Pager.getRow(i)
			if err != nil {
				return fmt.Errorf("failed to getRow(%d), %w", i, err)
			}
			fmt.Println(row)
		}
	}
	return nil
}
