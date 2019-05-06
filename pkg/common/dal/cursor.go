package dal

import (
	"encoding/json"
	"time"
)

type Cursor int64

func NewCursor(ts time.Time) Cursor {
	return Cursor(ts.UnixNano() / int64(time.Millisecond))
}

func (c Cursor) IsEmpty() bool {
	return c <= 0
}

func (c Cursor) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(c))
}

func (c Cursor) String() string {
	ts := time.Unix(int64(c), 0)

	return ts.Format(time.RFC1123)
}