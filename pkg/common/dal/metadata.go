package dal

import "time"

type (
	Entity struct {
		Metadata
		ID  string `json:"id"`
		Rev string `json:"rev"`
	}

	Metadata struct {
		CreatedAt time.Time `json:"created_at"`
		UpdateAt  time.Time `json:"update_at"`
	}
)

func (e Entity) IsEmpty() bool {
	return e.ID == ""
}
