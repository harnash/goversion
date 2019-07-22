package main

import (
	"os"

	"github.com/harnash/goversion/pkg/config"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app     = kingpin.New("goversion", "semantic versioning tool").Version("0.0.1")
	verbose = app.Flag("verbose", "Enable verbose mode").Bool()

	patch = app.Command("patch", "increments patch component of semantic version")
	minor = app.Command("minor", "increments minor component of semantic version")
	major = app.Command("major", "increments major component of semantic version")
)

func main() {
	cfg, _ := config.NewFromEnv()

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	cfg.VerboseMode = *verbose

	switch cmd {
	case patch.FullCommand():
		println("patch")
	case minor.FullCommand():
		println("minor")
	case major.FullCommand():
		println("major")
	}
	if cfg.VerboseMode {
		println("verbose mode on")
	}
}
