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
		cur := NewCursor(t, true)
		cur.setValue(s.RowToInsert)
		t.NumRows++

		return nil
	case StatementSelect:
		cur := NewCursor(t, false)
		for !cur.endOfTable {
			r, err := cur.getValue()
			if err != nil {
				return fmt.Errorf("failed to get row, %w", err)
			}
			fmt.Println(r)
			cur.advance()
		}
	}
	return nil
}
