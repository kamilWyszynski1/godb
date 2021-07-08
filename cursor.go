package godb

type Cursor struct {
	table  *Table
	rowNum uint32
	// Indicates a position one past the last element
	endOfTable bool
}

// NewCursor creates new Cursor
// if endFlag is set in create 'EndCursor' otherwise it creates 'StartCursor'
func NewCursor(table *Table, endFlag bool) *Cursor {
	c := &Cursor{
		table: table,
	}

	if endFlag {
		c.rowNum = table.NumRows
		c.endOfTable = true
	} else {
		c.rowNum = 0
		c.endOfTable = table.NumRows == 0
	}
	return c
}

func (c *Cursor) setValue(r *Row) {
	c.table.Pager.pages[c.rowNum] = r // write to wanted place
}

func (c *Cursor) advance() {
	c.rowNum++
	if c.rowNum == c.table.NumRows {
		c.endOfTable = true
	}
}

func (c *Cursor) getValue() (*Row, error) {
	return c.table.Pager.getRow(c.rowNum)
}
