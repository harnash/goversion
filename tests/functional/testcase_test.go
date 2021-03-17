package functional

import (
	"github.com/stretchr/testify/suite"
	"os"
	"path"
	"testing"
)

type TestCaseTest struct {
	BaseTestCase
}

func (suite *TestCaseTest) TestSimpleSuiteSetup() {
	// by now setup should be done already
	suite.NotEmpty(suite.baseFixturePath)
	testFolder := suite.baseFixturePath + "/TestSimpleSuiteSetup"
	suite.NotEmpty(testFolder)
	suite.NotEmpty(suite.fixturePath)
	suite.FolderMatchesDesiredState("fixtures were not setup properly")
}

func (suite *TestCaseTest) TestFileModifications() {
	filePath := path.Join(suite.FixturePath(), "someFile.txt")
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0660)
	suite.NoError(err, "could not open a test file for writing: %s", filePath)
	cnt, err := file.WriteString("\nNew line here!")
	suite.NoError(err, "could not write a test file: %s", filePath)
	suite.Equal(15, cnt, "wrote unexpected amount of bytes to test file: %s", filePath)
	suite.FolderMatchesDesiredState("file modifications were not detected properly")
}

func TestCaseTestSuite(t *testing.T) {
	testCase := new(TestCaseTest)
	curPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current working dir: %v", err)
	}
	testCase.SetBaseFixturePath(path.Join(curPath, "fixtures"))

	suite.Run(t, testCase)
}
