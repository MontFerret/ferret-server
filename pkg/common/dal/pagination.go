package dal

import "fmt"

type Pagination struct {
	Page uint `json:"page"`
	Size uint `json:"size"`
}

func (p *Pagination) String() string {
	return fmt.Sprintf("page %d size %d", p.Page, p.Size)
}
