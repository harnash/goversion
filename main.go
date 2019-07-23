package main

import (
	"os"
	"strings"

	"github.com/harnash/goversion/pkg/config"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	cfg, _ := config.NewFromEnv()

	app := kingpin.New("goversion", "semantic versioning tool").Version("0.0.1")
	app.Flag("verbose", "enable verbose mode").BoolVar(&cfg.VerboseMode)

	var noTag bool
	var noCommit bool

	app.Flag("dry-run", "enable dry-run mode (do not apply changes, use with 'verbose' mode").
		BoolVar(&cfg.DryRunMode)
	app.Flag("tag", "version will be tagged if possible (git/mercurial)").
		BoolVar(&cfg.TagVersion)
	app.Flag("no-tag", "should not create tag for new version").
		BoolVar(&noTag)
	app.Flag("allow-dirty", "proceed with bumping version even if repo has uncommitted changes").
		BoolVar(&cfg.AllowDirty)
	app.Flag("parse", "regex used to parse version string").StringVar(&cfg.ParseTemplate)
	app.Flag("serialize", "format used to print version string").StringsVar(&cfg.SerializeTemplate)
	app.Flag("commit", "should commit changes when finish").BoolVar(&cfg.CommitVersion)
	app.Flag("no-commit", "should not commit changes when finish").BoolVar(&noCommit)
	app.Flag("message", "template for message of the commit").StringVar(&cfg.CommitMessage)
	app.Flag("new-version", "new version that should be set").StringVar(&cfg.NewVersion)

	parts := app.Arg("part", "part of the version that should be bumped").Enum("patch", "minor", "major")
	files := app.Arg("files", "files to process").ExistingFiles()

	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	println(*parts)
	println(strings.Join(*files, "\n"))

	if cfg.VerboseMode {
		println("verbose mode on")
	}
	if cfg.TagVersion {
		println("tag version")
	}
}
