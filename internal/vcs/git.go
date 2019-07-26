package vcs

import "os/exec"

type GitVCS struct{}

func (GitVCS) String() string {
	return "GitVCS"
}

func (GitVCS) IsSupported() bool {
	cmd := exec.Command("git", "status")
	err := cmd.Run()
	if err != nil {
		return false
	}
	return cmd.ProcessState.Success()
}

func (GitVCS) IsDirty() bool {
	panic("implement me")
}

func (GitVCS) DoCommit(files []string, message string) error {
	panic("implement me")
}

func (GitVCS) MakeTag(name, message string) error {
	panic("implement me")
}

func init() {
	supportedVCS = append(supportedVCS, &GitVCS{})
}
