package db

type Settings struct {
	Endpoints []string
}

func NewDefaultSettings() Settings {
	return Settings{
		Endpoints: []string{
			"http://127.0.0.1:8529",
		},
	}
}
