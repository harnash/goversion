package config

import (
	"github.com/joomcode/errorx"
	"os"
	"regexp"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	CurrentVersion    string                 `yaml:"current_version"`
	NewVersion        string                 `yaml:"new_version"`
	CommitVersion     bool                   `yaml:"commit"`
	CommitMessage     string                 `yaml:"commit_message"`
	TagVersion        bool                   `yaml:"tag"`
	TagName           string                 `yaml:"tag_name"`
	VerboseMode       bool                   `yaml:"verbose"`
	DryRunMode        bool                   `yaml:"dry_run"`
	AllowDirty        bool                   `yaml:"allow_dirty"`
	List              bool                   `yaml:"list"`
	SerializeTemplate []string               `yaml:"serialize,flow"`
	ParseTemplate     *regexp.Regexp         `yaml:"parse"`
	ReleaseParts      map[string]ReleasePart `yaml:"parts,flow"`
	ReleaseFiles      map[string]ReleaseFile `yaml:"files,flow"`
}

type ReleasePart struct {
	OptionalValue string   `yaml:"optional_value,omitempty"`
	Values        []string `yaml:"values,flow"`
	FirstValue    string   `yaml:"first_value"`
}

type ReleaseFile struct {
	Search            string         `yaml:"search"`
	Replace           string         `yaml:"replace"`
	ParseTemplate     *regexp.Regexp `yaml:"parse"`
	SerializeTemplate []string       `yaml:"serialize,flow"`
}

func NewFromEnv() (*Configuration, error) {
	cfg := &Configuration{}
	return cfg, nil
}

func (c *Configuration) MergeWith(newConfig Configuration) error {
	return mergo.Merge(c, newConfig)
}

func (c Configuration) SaveToFile(file *os.File) error {
	bytes, err := yaml.Marshal(&c)
	if err != nil {
		return errorx.Decorate(err, "could not serialize config file")
	}

	_, err = file.Write(bytes)

	if err != nil {
		return errorx.Decorate(err, "could not save config to a file")
	}

	return nil
}

func NewFromFile(file *os.File) (*Configuration, error) {
	fileInfo, _ := (*file).Stat()
	data := make([]byte, fileInfo.Size())
	bytesRead, err := file.Read(data)
	if err != nil {
		return nil, err
	}

	fileConfig := &Configuration{}
	if err = yaml.Unmarshal(data[:bytesRead], fileConfig); err != nil {
		return nil, err
	}

	return fileConfig, err
}
