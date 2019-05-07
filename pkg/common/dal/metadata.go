package dal

import "time"

type Metadata struct {
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdateAt  *time.Time `json:"update_at,omitempty"`
}
