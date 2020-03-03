package main

import (
	"fmt"
	"github.com/harnash/goversion/pkg/bump"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"

	vcs2 "github.com/harnash/goversion/internal/vcs"

	"github.com/harnash/goversion/pkg/config"
)

func main() {
	cfg, _ := config.NewFromEnv()

	app := kingpin.New("goversion", "semantic versioning tool").Version("0.0.1")
	app.Flag("verbose", "enable verbose mode").BoolVar(&cfg.VerboseMode)

	var noTag bool
	var noCommit bool

	configFile := app.Flag("config", "file which contains configuration").Default(".goversion").File()
	app.Flag("dry-run", "enable dry-run mode (do not apply changes, use with 'verbose' mode").
		BoolVar(&cfg.DryRunMode)
	app.Flag("tag", "version will be tagged if possible (git/mercurial)").
		BoolVar(&cfg.TagVersion)
	app.Flag("tag-name", "sets the name of the tag for the new version").
		Default("v{{new_version}}").
		StringVar(&cfg.TagName)
	app.Flag("no-tag", "should not create tag for new version").
		BoolVar(&noTag)
	app.Flag("allow-dirty", "proceed with bumping version even if repo has uncommitted changes").
		BoolVar(&cfg.AllowDirty)
	app.Flag("parse", "regex used to parse version string").
		Default(`(?P<major>\d+)\.(?P<minor>\d+)\.(?P<patch>\d+)`).
		RegexpVar(&cfg.ParseTemplate)
	app.Flag("serialize", "format used to print version string").
		Default("{{.major}}.{{.minor}}.{{.patch}}").
		StringsVar(&cfg.SerializeTemplate)
	app.Flag("commit", "should commit changes when finish").BoolVar(&cfg.CommitVersion)
	app.Flag("no-commit", "should not commit changes when finish").BoolVar(&noCommit)
	app.Flag("message", "template for message of the commit").
		Default("Bump version: {{.current_version}} → {{.new_version}}").
		StringVar(&cfg.CommitMessage)
	app.Flag("new-version", "new version that should be set").StringVar(&cfg.NewVersion)

	part := app.Arg("part", "part of the version that should be bumped").Enum("patch", "minor", "major")
	files := app.Arg("files", "files to process").ExistingFiles()

	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	if configFile != nil {
		fileConfig, err := config.NewFromFile(*configFile)
		if err == nil {
			err = cfg.MergeWith(*fileConfig)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	vcs, err := vcs2.NewVCS()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(vcs)

	//TODO: Fetch current version parts
	version := bump.NewVersion([]bump.VersionPart{}, cfg.ParseTemplate, cfg.SerializeTemplate)
	err = version.Bump(*part)
	if err != nil {
		log.Fatal(err)
	}

	// add self config file to bump current version
	if configFile != nil {
		cfg.ReleaseFiles[(*configFile).Name()] = config.ReleaseFile{}
	}

	var filesToProcess []string
	if *files == nil {
		filesToProcess = make([]string, 0, len(cfg.ReleaseFiles))
		for key := range cfg.ReleaseFiles {
			filesToProcess = append(filesToProcess, key)
		}

	} else {
		filesToProcess = *files
	}

	err = bump.ApplyVersionToFiles(filesToProcess, version, cfg)
	if err != nil {
		log.Fatal(err)
	}
}
