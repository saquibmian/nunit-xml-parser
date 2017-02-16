package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"
)

type Test struct {
	Name            string
	Result          string
	TimeStarted     time.Time
	TimeEnded       time.Time
	DurationSeconds float64
	ErrorMessage    *string
	StackTrace      *string
}

func main() {
	var filePath string
	flag.StringVar(&filePath, "file", "", "the path to the file to parse")
	flag.Parse()
	if filePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	parsedXML, err := readXmlResults(filePath)
	if err != nil {
		panic(err)
	}

	tests, err := extractTests(parsedXML)
	if err != nil {
		panic(err)
	}
	err = json.NewEncoder(os.Stdout).Encode(tests)
	if err != nil {
		panic(err)
	}
}

func extractTests(results xmlTestResults) ([]Test, error) {
	timeStarted, err := results.TimeStarted()
	if err != nil {
		return nil, err
	}

	var tests []Test

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

func extractTestsRecursive(suite xmlTestSuite) []Test {
	var tests []Test

	for _, test := range suite.TestCases {
		tests = append(tests, Test{
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
