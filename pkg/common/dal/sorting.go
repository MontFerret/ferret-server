package dal

type (
	Order = uint64

	SortingField struct {
		Name  string `json:"name"`
		Order Order  `json:"order"`
	}

	Sorting struct {
		Fields []*SortingField `json:"columns"`
	}
)

const (
	OrderAsc  Order = 0
	OrderDesc Order = 1
)

var (
	OrderValue = map[Order]string{
		OrderAsc:  "ASC",
		OrderDesc: "DESC",
	}
)
