package main

import (
	"encoding/json"
	"flag"
	"os"

	"./model"
	_ "./nunit"
	_ "./jenkins"
)

func main() {
	var filePath string
	var format string
	flag.StringVar(&format, "format", "", "the format of the file to parse")
	flag.StringVar(&filePath, "file", "", "the path to the file to parse")
	flag.Parse()
	if format == "" || filePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	tests, err := model.Extract(format, filePath)
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(os.Stdout).Encode(tests)
	if err != nil {
		panic(err)
	}
}
