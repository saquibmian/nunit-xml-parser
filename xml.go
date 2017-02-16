package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"time"
)

const (
	DateTimeFormat = "2006-01-02T15:04:05"
)

type XmlTestResults struct {
	Total        int `xml:"total,attr"`
	Errors       int `xml:"errors,attr"`
	Failures     int `xml:"failures,attr"`
	Ignored      int `xml:"ignored,attr"`
	Inconclusive int `xml:"inconclusive,attr"`
	Skipped      int `xml:"skipped,attr"`

	DateStartedString string `xml:"date,attr"`
	TimeStartedString string `xml:"time,attr"`

	TestSuites []XmlTestSuite `xml:"test-suite"`
}

func (r *XmlTestResults) TimeStarted() (time.Time, error) {
	timeStarted, err := time.Parse(DateTimeFormat, r.DateStartedString+"T"+r.TimeStartedString)
	if err != nil {
		return timeStarted, fmt.Errorf("error parsing time: %s", err)
	}
	timeStarted = timeStarted.UTC()
	return timeStarted, nil
}

type XmlTestSuite struct {
	Type   string  `xml:"type,attr"`
	Name   string  `xml:"name,attr"`
	Result string  `xml:"result,attr"`
	Time   float64 `xml:"time,attr"`

	TestSuites []XmlTestSuite `xml:"results>test-suite"`
	TestCases  []XmlTestCase  `xml:"results>test-case"`
}

type XmlTestCase struct {
	Name           string         `xml:"name,attr"`
	Result         string         `xml:"result,attr"`
	Time           float64        `xml:"time,attr"`
	FailureMessage XmlNestedCData `xml:"failure>message"`
	StackTrace     XmlNestedCData `xml:"failure>stack-trace"`
}

type XmlNestedCData struct {
	Contents *string `xml:",chardata"`
}

func readXmlResults(filePath string) (XmlTestResults, error) {
	results := XmlTestResults{}

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
