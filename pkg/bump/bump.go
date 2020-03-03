package bump

import (
	"bytes"
	"fmt"
	"github.com/joomcode/errorx"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"text/template"

	"github.com/harnash/goversion/pkg/config"
)

type VersionPart struct {
	name string
	value string
	customPart config.ReleasePart //TODO: move out from config
}

func NewVersionPart(name, value string, customPart config.ReleasePart) *VersionPart {
	return &VersionPart{
		name,
		value,
		customPart,
	}
}

type Version struct {
	parts []VersionPart
	parseTemplate *regexp.Regexp
	serializeTemplate []string
}
// TODO: serialization to proper string, serialization of current version, proper major bumping - reset minor etc.
// TODO: "constructor"
// TODO: Test, test, tests

func NewVersion(parts []VersionPart, parseTemplate *regexp.Regexp, serializeTemplate []string) *Version {
	return &Version{
		parts,
		parseTemplate,
		serializeTemplate,
	}
}

func (v Version) Serialize() ([]byte, error) {
	buff := bytes.NewBufferString("")
	for idx, serializeTemplate := range v.serializeTemplate {
		tmpl, err := template.New(fmt.Sprintf("template_%d", idx)).Parse(serializeTemplate)
		if err != nil {
			return nil, errorx.Decorate(err, "invalid version serialization template")
		}

		var data = make(map[string]string, len(v.parts));
		for _, part := range v.parts {
			data[part.name] = part.value
		}
		err = tmpl.Execute(buff, data)
		if err != nil {
			continue
		}
		break
	}

	return buff.Bytes(), nil
}

func (v *Version) Bump(part string) error {
	serializedVersion, err := v.Serialize()
	if err != nil {
		return errorx.Decorate(err, "could not serialize version")
	}
	matchedVersions := v.parseTemplate.FindStringSubmatch(string(serializedVersion))
	if matchedVersions == nil {
		return errorx.IllegalArgument.New("could not parse current version")
	}

	if len(v.parseTemplate.SubexpNames()) == 0 {
		return errorx.IllegalFormat.New("parse template has no meaningful part names")
	}

	var partValue string
	versionParts := map[string]string{}
	for idx, partName := range v.parseTemplate.SubexpNames() {
		if partName == "" {
			continue
		}
		versionParts[partName] = matchedVersions[idx]
		if partName == part {
			partValue = matchedVersions[idx]
		}
	}

	if partValue == "" {
		return errorx.IllegalArgument.New("could not find version part to bump")
	}

	found := false
	for _, partConfig := range v.parts {
		if partConfig.name == part {
			found = true
		}

		if !found {
			continue
		}

		nextValue := partConfig.customPart.FirstValue
		foundValue := false
		for _, value := range partConfig.customPart.Values {
			if value == partValue {
				foundValue = true
			} else if foundValue {
				nextValue = value
				break
			}
		}
		versionParts[part] = nextValue
	}

	if !found {
		intVersion, err := strconv.Atoi(partValue)
		if err != nil {
			return errorx.IllegalFormat.New("could not parse version part to integer")
		}
		versionParts[part] = strconv.Itoa(intVersion+1)

	}

	return nil
}

func (v Version) ApplyFileConfiguration(fileName string, configuration *config.Configuration) Version {
	releaseFile, ok := configuration.ReleaseFiles[fileName]
	if !ok {
		return v
	}

	newVersion := v
	if len(releaseFile.SerializeTemplate) != 0 {
		newVersion.serializeTemplate = releaseFile.SerializeTemplate
	}

	if releaseFile.ParseTemplate != nil {
		newVersion.parseTemplate = releaseFile.ParseTemplate
	}

	return newVersion
}

func ApplyVersionToFiles(files []string, newVersion *Version, configuration *config.Configuration) error {
	for _, file := range files {
		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return errorx.Decorate(err, fmt.Sprintf("could not open file: %s", file))
		}

		fileInfo, err := os.Stat(file)
		if err != nil {
			return errorx.Decorate(err, fmt.Sprintf("could not get stat info for file: %s", file))
		}

		ver := newVersion.ApplyFileConfiguration(file, configuration)

		versionSerialized, err := ver.Serialize()
		if err != nil {
			return errorx.Decorate(err, "could not serialize new version")
		}

		if len(versionSerialized) == 0 {
			return errorx.IllegalFormat.New("could not serialize new version using any available serialization templates")
		}

		contents = ver.parseTemplate.ReplaceAllLiteral(contents, versionSerialized)
		err = ioutil.WriteFile(file, contents, fileInfo.Mode())
		if err != nil {
			return errorx.Decorate(err, fmt.Sprintf("could not write file: %s", file))
		}

	}

	return nil
}
