package functional

import (
	"github.com/joomcode/errorx"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type TestCaseTest struct {
	BaseTestCase
}

func (suite *TestCaseTest) compareDir(src, dst string) error {
	srcContents, err := ioutil.ReadDir(src)
	if err != nil {
		return errorx.Decorate(err, "could not read contents of source directory: %s", src)
	}

	for _, srcItem := range srcContents {
		srcPath := path.Join(src, srcItem.Name())
		dstPath := path.Join(dst, srcItem.Name())
		if srcItem.IsDir() {
			err := suite.compareDir(srcPath, dstPath)
			if err != nil {
				return errorx.Decorate(err, "source and destination paths are different: %s, %s", srcPath, dstPath)
			}
		}

		dstItem, err := os.Stat(dstPath)

		if suite.NoError(err, "could not find file in destination folder: %s", dstPath) {
			suite.Equal(srcItem.Size(), dstItem.Size(), "sizes of source and destination files differ: %s, %s", srcPath, dstPath)
		}
	}

	return nil
}

func (suite *TestCaseTest) TestSimpleSuiteSetup() {
	// by now setup should be done already
	suite.NotEmpty(suite.sourceFixturePath)
	suite.NotEmpty(suite.fixturePath)

	err := suite.compareDir(suite.sourceFixturePath, suite.fixturePath)
	suite.NoError(err)
}


func TestCaseTestSuite(t *testing.T) {
	testCase := new(TestCaseTest)
	curPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current working dir: %v", err)
	}
	testCase.SetFixture(path.Join(curPath, "fixtures/TestCase01"))

	suite.Run(t, testCase)
}
