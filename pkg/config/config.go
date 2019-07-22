package config

type Configuration struct {
	CurrentVersion string
	NewVersion     string
	CommitVersion  bool
	CommitMessage  string
	TagVersion     bool
	TagName        string
	VerboseMode    bool
}

func NewFromEnv() (*Configuration, error) {
	cfg := &Configuration{}
	return cfg, nil
}
