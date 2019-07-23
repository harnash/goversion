package config

import (
	"os"

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
	SerializeTemplate []string               `yaml:"serialize"`
	ParseTemplate     string                 `yaml:"parse"`
	ReleaseParts      map[string]ReleasePart `yaml:"parts"`
	ReleaseFiles      map[string]ReleaseFile `yaml:"files"`
}

type ReleasePart struct {
	OptionalValue string   `yaml:"optional_value"`
	Values        []string `yaml:"values"`
	FirstValue    string   `yaml:"first_value"`
}

type ReleaseFile struct {
	Search            string   `yaml:"search"`
	Replace           string   `yaml:"replace"`
	ParseTemplate     string   `yaml:"parse"`
	SerializeTemplate []string `yaml:"serialize"`
}

func NewFromEnv() (*Configuration, error) {
	cfg := &Configuration{}
	return cfg, nil
}

func (c *Configuration) MergeWith(newConfig Configuration) error {
	return mergo.Merge(c, newConfig)
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
