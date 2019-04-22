package dal

type Entity struct {
	Metadata
	ID  string `json:"id"`
	Rev string `json:"rev"`
}

func (e Entity) IsEmpty() bool {
	return e.ID == ""
}
