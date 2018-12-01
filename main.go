package main

import (
	"fmt"
	"github.com/MontFerret/ferret-server/server"
	"github.com/namsral/flag"
	"os"
	"strings"
)

var defaultSettings = server.NewDefaultSettings()

type Params []string

func (p *Params) String() string {
	return "[" + strings.Join(*p, ",") + "]"
}

func (p *Params) Set(value string) error {
	*p = append(*p, value)
	return nil
}

var (
	version = "undefined"

	showHelp = flag.Bool(
		"help",
		false,
		"show this list",
	)

	showVersion = flag.Bool(
		"version",
		false,
		"show version",
	)

	// General
	name = flag.String(
		"name",
		defaultSettings.Name,
		"server name",
	)

	cdpAddress = flag.String(
		"cdp",
		defaultSettings.CDPAddress,
		"CDP address",
	)

	// HTTP
	httpPort = flag.Uint64(
		"http-port",
		defaultSettings.HTTP.Port,
		"http server port number",
	)

	// Db
	dbEndpoints Params
)

func createSettings() (server.Settings, error) {
	settings := server.NewDefaultSettings()

	// General
	settings.Name = *name
	settings.Version = version
	settings.CDPAddress = *cdpAddress

	// HTTP
	settings.HTTP.Port = *httpPort

	// Db
	if len(dbEndpoints) > 0 {
		settings.Database.Endpoints = dbEndpoints
	}

	return settings, nil
}

func main() {
	flag.Var(
		&dbEndpoints,
		"db",
		`database endpoint (--db=http://0.0.0.0:8529, --db=http://0.0.0.0:8530)`,
	)

	flag.Parse()

	if *showHelp {
		flag.PrintDefaults()
		os.Exit(0)
		return
	}

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
		return
	}

	settings, err := createSettings()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	a, err := server.New(settings)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}
}
