package dal

import "fmt"

type Pagination struct {
	Cursor Cursor `json:"cursor"`
	Count  uint64 `json:"count"`
}

func (p *Pagination) String() string {
	return fmt.Sprintf("count %d cursor %s", p.Count, p.Cursor)
}
