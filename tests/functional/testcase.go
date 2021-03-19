package functional

import (
	"fmt"
	"github.com/joomcode/errorx"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"path"
)

const inputFolderName = "in"
const outputFolderName = "out"

type BaseTestCase struct {
	suite.Suite
	baseFixturePath   string
	fixturePath       string
	desiredOutputPath string
}

func (tc *BaseTestCase) FixturePath() string {
	return tc.fixturePath
}

//SetBaseFixturePath will set path to the source input files which will be modified by the test
func (tc *BaseTestCase) SetBaseFixturePath(path string) {
	tc.baseFixturePath = path
}

//BaseFixturePath return path to the source input files
func (tc *BaseTestCase) BaseFixturePath() string {
	return tc.baseFixturePath
}

func (tc *BaseTestCase) BeforeTest(suiteName, testName string) {
	var err error
	tc.fixturePath, err = ioutil.TempDir("", fmt.Sprintf("fixtures-%s-*", suiteName))
	tc.NoError(err, "Could not create temporary directory for fixtures")

	testFolder := path.Join(tc.BaseFixturePath(), testName, inputFolderName)
	if !tc.DirExists(testFolder, "source fixture folder '%s' does not exists", testFolder) {
		tc.FailNow("fixtures not configured properly")
	}

	err = tc.copyDir(testFolder, tc.fixturePath)
	if !tc.NoError(err, "Could not copy fixtures from %s", testFolder) {
		tc.FailNow("failed to setup test suite fixtures")
	}

	tc.desiredOutputPath = path.Join(tc.BaseFixturePath(), testName, outputFolderName)
}

func (tc BaseTestCase) AfterTest(suiteName, testName string) {
	if tc.fixturePath != "" {
		err := os.RemoveAll(tc.fixturePath)
		tc.NoError(err, "Could not remove temporary fixture directory: %s", tc.fixturePath)
	}
}

func (tc *BaseTestCase) compareDir(src, dst string) error {
	srcContents, err := ioutil.ReadDir(src)
	if err != nil {
		return errorx.Decorate(err, "could not read contents of source directory: %s", src)
	}

	for _, srcItem := range srcContents {
		srcPath := path.Join(src, srcItem.Name())
		dstPath := path.Join(dst, srcItem.Name())
		if srcItem.IsDir() {
			err := tc.compareDir(srcPath, dstPath)
			if err != nil {
				return errorx.Decorate(err, "source and destination paths are different: %s, %s", srcPath, dstPath)
			}
		}

		dmp := diffmatchpatch.New()
		srcData, err := ioutil.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("could not read contents of a file: %v: %v", srcPath, err)
		}
		dstData, err := ioutil.ReadFile(dstPath)
		if err != nil {
			return fmt.Errorf("could not read contents of a file: %v: %v", dstPath, err)
		}

		diffData := dmp.DiffMain(string(dstData), string(srcData), true)
		for _, diff := range diffData {
			if diff.Type != diffmatchpatch.DiffEqual {
				return fmt.Errorf("files differ: %s != %s:\n%s",
					srcPath, dstPath, dmp.DiffPrettyText(diffData))
			}
		}
	}

	return nil
}

func (tc BaseTestCase) FolderMatchesDesiredState(msgAndArgs ...interface{}) bool {
	err := tc.compareDir(tc.fixturePath, tc.desiredOutputPath)
	return tc.NoError(err, msgAndArgs...)
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
