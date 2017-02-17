package jenkins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/saquibmian/tests-parser/model"
)

func init() {
	model.RegisterExtractor("jenkins", jenkinsExtractor{})
}

const (
	dateTimeFormat = "2006-01-02T15:04:05"
)

type jenkinsRoot struct {
	Suites []jenkinsSuite `json:"suites"`
}

type jenkinsSuite struct {
	Cases []jenkinsCase `json:"cases"`
}

type jenkinsCase struct {
	Name            string  `json:"name"`
	FixtureName     string  `json:"className"`
	Result          string  `json:"status"`
	Message         string  `json:"stderr"`
	DurationSeconds float64 `json:"duration"`
	ErrorMessage    *string `json:"errorDetails"`
	StackTrace      *string `json:"errorStackTrace"`
}

type jenkinsExtractor struct{}

// ParseXMLResults converts an NUnit XML results file into the test model
func (e jenkinsExtractor) Extract(filePath string) ([]model.Test, error) {
	var tests []model.Test

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return tests, fmt.Errorf("unable to read file: %s", err)
	}

	results := jenkinsRoot{}
	err = json.Unmarshal(contents, &results)
	if err != nil {
		return tests, fmt.Errorf("unable to parse file: %s", err)
	}

	tests, err = extractTests(results)
	if err != nil {
		return tests, fmt.Errorf("unable to extract test results: %s", err)
	}

	return tests, nil
}

func extractTests(results jenkinsRoot) ([]model.Test, error) {
	var tests []model.Test

	for _, suite := range results.Suites {
		for _, testCase := range suite.Cases {
			tests = append(tests, model.Test{
				Name:            testCase.FixtureName + "." + testCase.Name,
				FixtureName:     testCase.FixtureName,
				Result:          testCase.Result,
				DurationSeconds: testCase.DurationSeconds,
				ErrorMessage:    testCase.ErrorMessage,
				StackTrace:      testCase.StackTrace,
			})
		}
	}

	return tests, nil
}
