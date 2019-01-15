package db

type Settings struct {
	Endpoints []string
}

func NewDefaultSettings() Settings {
	return Settings{
		Endpoints: []string{
			"http://127.0.0.1:8529",
			"http://0.0.0.0:8529",
			"http://arangodb:8529",
		},
	}
}
