package vcs

import "os/exec"

type MercurialVCS struct{}

func (MercurialVCS) String() string {
	return "MercurialVCS"
}

func (MercurialVCS) IsSupported() bool {
	cmd := exec.Command("hg", "status")
	err := cmd.Run()
	if err != nil {
		return false
	}

	return cmd.ProcessState.Success()
}

func (MercurialVCS) IsDirty() bool {
	panic("implement me")
}

func (MercurialVCS) DoCommit(files []string, message string) error {
	panic("implement me")
}

func (MercurialVCS) MakeTag(name, message string) error {
	panic("implement me")
}

func init() {
	supportedVCS = append(supportedVCS, &MercurialVCS{})
}
