package dal

type QueryResult struct {
	BeforeCursor Cursor `json:"before_cursor"`
	AfterCursor  Cursor `json:"after_cursor"`
	Count        uint64 `json:"count"`
}
