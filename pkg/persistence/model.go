package persistence

type (
	Local string

	Remote struct {
		Kind   string                 `json:"kind"`
		Params map[string]interface{} `json:"params"`
	}
)
