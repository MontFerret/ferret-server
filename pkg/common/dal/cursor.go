package dal

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
)

type Cursor string

func (c Cursor) IsEmpty() bool {
	return c == ""
}

func (c Cursor) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(c))
}

func (c Cursor) String() string {
	return string(c)
}

func EncodeCursor(value int64) Cursor {
	str := strconv.FormatInt(value, 10)
	encoded := base64.StdEncoding.EncodeToString([]byte(str))

	return Cursor(encoded)
}

func DecodeCursor(c Cursor) int64 {
	if c == "" {
		return 0
	}

	decoded, err := base64.StdEncoding.DecodeString(string(c))

	if err != nil {
		return 0
	}

	num, err := strconv.ParseInt(string(decoded), 10, 64)

	if err != nil {
		return 0
	}

	return num
}
