package bump

import (
	"bytes"
	"fmt"
	"github.com/joomcode/errorx"
	"io/ioutil"
	"os"
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
		if partName == "" {
			continue
		}
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
		versionParts[part] = strconv.Itoa(intVersion+1)

	}

	return versionParts, nil
}

func SerializeVersion(versionParts map[string]string, serializeTemplate []string) ([]byte, error) {
	buff := bytes.NewBufferString("")
	for idx, serializeTemplate := range serializeTemplate {
		tmpl, err := template.New(fmt.Sprintf("template_%d", idx)).Parse(serializeTemplate)
		if err != nil {
			return nil, errorx.Decorate(err, "invalid version serialization template")
		}

		err = tmpl.Execute(buff, versionParts)
		if err != nil {
			continue
		}
		break
	}

	return buff.Bytes(), nil
}

func ApplyVersionToFiles(files []string, newVersionParts map[string]string, configuration *config.Configuration) error {
	for _, file := range files {
		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return errorx.Decorate(err, fmt.Sprintf("could not open file: %s", file))
		}

		fileInfo, err := os.Stat(file)
		if err != nil {
			return errorx.Decorate(err, fmt.Sprintf("could not get stat info for file: %s", file))
		}

		releaseFile, ok := configuration.ReleaseFiles[file]
		if !ok {
			releaseFile = config.ReleaseFile{
				ParseTemplate:     configuration.ParseTemplate,
				SerializeTemplate: configuration.SerializeTemplate,
			}
		}
		if len(releaseFile.SerializeTemplate) == 0 {
			releaseFile.SerializeTemplate = configuration.SerializeTemplate
		}
		if releaseFile.ParseTemplate == nil {
			releaseFile.ParseTemplate = configuration.ParseTemplate
		}

		versionSerialized, err := SerializeVersion(newVersionParts, releaseFile.SerializeTemplate)
		if err != nil {
			return errorx.Decorate(err, "could not serialize new version")
		}

		if len(versionSerialized) == 0 {
			return errorx.IllegalFormat.New("could not serialize new version using any available serialization templates")
		}

		contents = releaseFile.ParseTemplate.ReplaceAllLiteral(contents, versionSerialized)
		err = ioutil.WriteFile(file, contents, fileInfo.Mode())
		if err != nil {
			return errorx.Decorate(err, fmt.Sprintf("could not write file: %s", file))
		}

	}

	return nil
}
