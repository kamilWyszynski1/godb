package godb

import "os"

type Pager struct {
	f          *os.File
	fileLength uint32
	pages      []*Row
}
