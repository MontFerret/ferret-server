package server

import (
	"github.com/MontFerret/ferret-server/server/db"
	"github.com/MontFerret/ferret-server/server/http"
)

type Settings struct {
	Version    string
	Name       string
	CDPAddress string
	HTTP       http.Settings
	Database   db.Settings
}

func NewDefaultSettings() Settings {
	return Settings{
		Name:       "ferret-server",
		CDPAddress: "http://0.0.0.0:9222",
		HTTP:       http.NewDefaultSettings(),
		Database:   db.NewDefaultSettings(),
	}
}
