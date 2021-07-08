package godb

import (
	"errors"
	"fmt"
	"os"
)

const PageSize = 4096

type Pager struct {
	f          *os.File
	fileLength uint32
	pages      [TableMaxRows]*Row
}

func (p *Pager) getRow(pageNum uint32) (*Row, error) {
	if pageNum > TableMaxRows {
		return nil, errors.New("pageNum exceeds TableMaxRows value")
	}

	if p.pages[pageNum] != nil {
		return p.pages[pageNum], nil
	} else {
		// read and parse Row from file
		read := make([]byte, RowSize())
		_, err := p.f.ReadAt(read, int64(pageNum*RowSize()))
		if err != nil {
			return nil, fmt.Errorf("failed to read wanted page from file, %w", err)
		}
		var r Row
		if err := UnmarshalRow(read, &r); err != nil {
			return nil, fmt.Errorf("failed to unmarshal row, %w", err)
		}
		p.pages[pageNum] = &r
		return &r, nil
	}
}

func (p *Pager) Close(pagesNum uint32) error {
	logger.
		WithField("method", "Close").
		Infof("closing pager, from: %d to: %d\n", p.fileLength/RowSize(), pagesNum)
	for i := p.fileLength / RowSize(); i < pagesNum; i++ {
		page := p.pages[i]
		data, err := page.Marshal()
		if err != nil {
			return fmt.Errorf("failed to marshal row, %w", err)
		}
		if _, err := p.f.Write(data); err != nil {
			return fmt.Errorf("failed to write to file, %w", err)
		}
	}
	return nil
}
