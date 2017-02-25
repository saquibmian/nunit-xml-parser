package nunit

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/saquibmian/tests-parser/model"
)

func init() {
	err := model.RegisterExtractor("nunit", nunitExtractor{})
	if err != nil {
		panic(err)
	}
}

const (
	dateTimeFormat = "2006-01-02T15:04:05"
)

type xmlTestResults struct {
	Total        int `xml:"total,attr"`
	Errors       int `xml:"errors,attr"`
	Failures     int `xml:"failures,attr"`
	Ignored      int `xml:"ignored,attr"`
	Inconclusive int `xml:"inconclusive,attr"`
	Skipped      int `xml:"skipped,attr"`

	DateStartedString string `xml:"date,attr"`
	TimeStartedString string `xml:"time,attr"`

	TestSuites []xmlTestSuite `xml:"test-suite"`
}

func (r *xmlTestResults) TimeStarted() (time.Time, error) {
	timeStarted, err := time.Parse(dateTimeFormat, r.DateStartedString+"T"+r.TimeStartedString)
	if err != nil {
		return timeStarted, fmt.Errorf("error parsing time: %s", err)
	}
	timeStarted = timeStarted.UTC()
	return timeStarted, nil
}

type xmlTestSuite struct {
	Type   string  `xml:"type,attr"`
	Name   string  `xml:"name,attr"`
	Result string  `xml:"result,attr"`
	Time   float64 `xml:"time,attr"`

	TestSuites []xmlTestSuite `xml:"results>test-suite"`
	TestCases  []xmlTestCase  `xml:"results>test-case"`
}

type xmlTestCase struct {
	Name           string         `xml:"name,attr"`
	Result         string         `xml:"result,attr"`
	Time           float64        `xml:"time,attr"`
	FailureMessage xmlNestedCData `xml:"failure>message"`
	StackTrace     xmlNestedCData `xml:"failure>stack-trace"`
}

type xmlNestedCData struct {
	Contents *string `xml:",chardata"`
}

type nunitExtractor struct{}

// ParseXMLResults converts an NUnit XML results file into the test model
func (e nunitExtractor) Extract(filePath string) ([]model.Test, error) {
	var tests []model.Test

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return tests, fmt.Errorf("unable to read file: %s", err)
	}

	results := xmlTestResults{}
	err = xml.Unmarshal(contents, &results)
	if err != nil {
		return tests, fmt.Errorf("unable to parse file: %s", err)
	}

	tests, err = extractTests(results)
	if err != nil {
		return tests, fmt.Errorf("unable to extract test results: %s", err)
	}

	return tests, nil
}

func extractTests(results xmlTestResults) ([]model.Test, error) {
	timeStarted, err := results.TimeStarted()
	if err != nil {
		return nil, err
	}

	var tests []model.Test

	for _, suite := range results.TestSuites {
		innerTests := extractTestsRecursive(suite)

		// adjust time started and ended
		for _, test := range innerTests {
			test.TimeStarted = timeStarted
			test.TimeEnded = timeStarted.Add(time.Duration(test.DurationSeconds) * time.Second)
			tests = append(tests, test)
		}
	}

	return tests, nil
}

func extractTestsRecursive(suite xmlTestSuite) []model.Test {
	var tests []model.Test

	for _, test := range suite.TestCases {
		tests = append(tests, model.Test{
			Name:            test.Name,
			Result:          test.Result,
			DurationSeconds: test.Time,
			ErrorMessage:    test.FailureMessage.Contents,
			StackTrace:      test.StackTrace.Contents,
		})
	}

	for _, innerSuite := range suite.TestSuites {
		innerTests := extractTestsRecursive(innerSuite)
		tests = append(tests, innerTests...)
	}

	return tests
}
