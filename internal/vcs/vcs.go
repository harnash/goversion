package vcs

import "fmt"

type VersionControlSystem interface {
	IsSupported() bool
	IsDirty() bool
	DoCommit(files []string, message string) error
	MakeTag(name, message string) error
	String() string
}

var supportedVCS []VersionControlSystem

func NewVCS() (VersionControlSystem, error) {
	for _, vcs := range supportedVCS {
		if vcs.IsSupported() {
			return vcs, nil
		}
	}

	return nil, fmt.Errorf("could not find supported VCS in the working directory")
}
