package bump

import (
	"github.com/harnash/goversion/pkg/config"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestVersion_Serialize(t *testing.T) {
	relPart := config.ReleasePart{
		OptionalValue: "",
		Values:        nil,
		FirstValue:    "",
	}
	major := NewVersionPart("major", "0", relPart)
	minor := NewVersionPart("minor", "0", relPart)
	patch := NewVersionPart("patch", "0", relPart)
	v := NewVersion(
		[]VersionPart{*major, *minor, *patch},
		regexp.MustCompile(""),
		[]string{"v{{.major}}.{{.minor}}.{{.patch}}"},
	)

	buff, err := v.Serialize()
	if assert.NoError(t, err) {
		assert.Equal(t, "v0.0.0", string(buff))
	}
}
