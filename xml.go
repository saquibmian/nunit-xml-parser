package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"time"
)

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

func readXmlResults(filePath string) (xmlTestResults, error) {
	results := xmlTestResults{}

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return results, fmt.Errorf("unable to read file: %s", err)
	}

	err = xml.Unmarshal(contents, &results)
	if err != nil {
		return results, fmt.Errorf("unable to parse file: %s", err)
	}

	return results, nil
}
