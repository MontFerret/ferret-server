package dal

import (
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type Cursor string

func NewCursor(ts time.Time) Cursor {
	str := strconv.Itoa(int(ts.Unix()))
	encoded := base64.StdEncoding.EncodeToString([]byte(str))

	return Cursor(encoded)
}

func (c Cursor) IsEmpty() bool {
	return c == ""
}

func (c Cursor) MarshalJSON() ([]byte, error) {
	if c == "" {
		return []byte(""), nil
	}

	ts, err := DecodeCursor(c)

	if err != nil {
		return nil, err
	}

	return json.Marshal(ts)
}

func (c Cursor) String() string {
	return string(c)
}

func DecodeCursor(c Cursor) (time.Time, error) {
	if c == "" {
		return time.Now(), nil
	}

	decoded, err := base64.StdEncoding.DecodeString(string(c))

	if err != nil {
		return time.Time{}, errors.Wrap(err, "decode cursor")
	}

	num, err := strconv.ParseInt(string(decoded), 10, 64)

	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(int64(num), -1), nil
}
