package functional

import (
	"github.com/joomcode/errorx"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"path"
)

type BaseTestCase struct {
	suite.Suite
	sourceFixturePath string
	fixturePath string
}


func (tc *BaseTestCase) SetFixture(path string) {
	tc.sourceFixturePath = path
}

func (tc *BaseTestCase) SetupSuite() {
	if tc.sourceFixturePath != "" {
		var err error

		tc.fixturePath, err = ioutil.TempDir("", "fixtures-*")
		tc.NoError(err, "Could not create temporary directory for fixtures")

		err = tc.copyDir(tc.sourceFixturePath, tc.fixturePath)
		tc.NoError(err, "Could not copy fixtures from %s", tc.sourceFixturePath)
	}
}

func (tc BaseTestCase) TearDownSuite() {
	if tc.fixturePath != "" {
		err := os.RemoveAll(tc.fixturePath)
		tc.NoError(err, "Could not remove temporary fixture directory: %s", tc.fixturePath)
	}
}

func (tc BaseTestCase) copyDir(src, dst string) error {
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return errorx.Decorate(err, "could not read contents of the source directory: %s", src)
	}

	err = os.MkdirAll(dst, 0777)
	if err != nil {
		return errorx.Decorate(err, "could not create destination directory: %s", dst)
	}

	for _, item := range entries {
		srcPath := path.Join(src, item.Name())
		dstPath := path.Join(dst, item.Name())

		if item.IsDir() {
			err = tc.copyDir(srcPath, dstPath)
			if err != nil {
				return errorx.Decorate(err, "could not copy directory contents: %s -> %s", srcPath, dstPath)
			}
			continue
		}

		contents, err := ioutil.ReadFile(srcPath)
		if err != nil {
			return errorx.Decorate(err, "could not read contents of a file: %s", srcPath)
		}

		err = ioutil.WriteFile(dstPath, contents, 0777)
		if err != nil {
			return errorx.Decorate(err, "could not write contents of a file: %s", dstPath)
		}
	}

	return nil
}
