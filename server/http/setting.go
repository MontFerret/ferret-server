package http

type Settings struct {
	Port uint64
}

func NewDefaultSettings() Settings {
	return Settings{
		Port: 8080,
	}
}
