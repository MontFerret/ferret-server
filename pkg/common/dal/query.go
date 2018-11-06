package dal

type Query struct {
	Pagination Pagination `json:"pagination"`
	Sorting    Sorting    `json:"sorting"`
	Filtering  Filtering  `json:"filtering"`
}
