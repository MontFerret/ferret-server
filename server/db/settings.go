package db

type Settings struct {
	Endpoints []string
}

func NewDefaultSettings() Settings {
	return Settings{}
}
