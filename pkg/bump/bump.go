package bump

import (
	"bytes"
	"fmt"
	"github.com/joomcode/errorx"
	"io/ioutil"
	"strconv"
	"text/template"

	"github.com/harnash/goversion/pkg/config"
)

func VersionBump(part string, configuration *config.Configuration) (map[string]string, error) {
	matchedVersions := configuration.ParseTemplate.FindStringSubmatch(configuration.CurrentVersion)
	if matchedVersions == nil {
		return nil, errorx.IllegalArgument.New("could not parse current version")
	}

	if len(configuration.ParseTemplate.SubexpNames()) == 0 {
		return nil, errorx.IllegalFormat.New("parse template has no meaningful part names")
	}

	var partValue string
	versionParts := map[string]string{}
	for idx, partName := range configuration.ParseTemplate.SubexpNames() {
		versionParts[partName] = matchedVersions[idx]
		if partName == part {
			partValue = matchedVersions[idx]
		}
	}

	if partValue == "" {
		return nil, errorx.IllegalArgument.New("could not find version part to bump")
	}

	partConfig, ok := configuration.ReleaseParts[part]
	if ok {
		nextValue := partConfig.FirstValue
		foundValue := false
		for _, value := range partConfig.Values {
			if value == partValue {
				foundValue = true
			} else if foundValue {
				nextValue = value
				break
			}
		}
		versionParts[part] = nextValue
	} else {
		intVersion, err := strconv.Atoi(partValue)
		if err != nil {
			return nil, errorx.IllegalFormat.New("could not parse version part to integer")
		}
		versionParts[part] = string(intVersion+1)

	}

	return versionParts, nil
}

func ApplyVersionToFiles(files []string, newVersionParts map[string]string, configuration config.Configuration) error {
	for _, file := range files {
		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return errorx.Decorate(err, fmt.Sprintf("could not open file: %s", file))
		}

		releaseFile, ok := configuration.ReleaseFiles[file]
		if ok {

		} else {
			releaseFile = config.ReleaseFile{
				ParseTemplate:     configuration.ParseTemplate,
				SerializeTemplate: configuration.SerializeTemplate,
			}
		}

		var newVersion string
		buff := bytes.NewBufferString(newVersion)
		for idx, serializeTemplate := range releaseFile.SerializeTemplate {
			tmpl, err := template.New(fmt.Sprintf("file_%s_temaplte_%d", file, idx)).Parse(serializeTemplate)
			if err != nil {
				return errorx.Decorate(err, "invalid version serialization template")
			}

			err = tmpl.Execute(buff, newVersionParts)
			if err != nil {
				continue
			}
		}

		if newVersion == "" {
			return errorx.IllegalFormat.New("could not serialize new version using any available serialization templates")
		}

		contents = releaseFile.ParseTemplate.ReplaceAllLiteral(contents, []byte(newVersion))
		err = ioutil.WriteFile(file, contents, 0644)
		if err != nil {
			return errorx.Decorate(err, fmt.Sprintf("could not write file: %s", file))
		}

	}

	return nil
}
