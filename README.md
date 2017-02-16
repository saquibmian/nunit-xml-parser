# nunit-xml-parser

A parser for the NUnit XML test results format. The tool dumps to JSON.

# usage

`-file {file}` will output an array of test results defined in the results file.

```
$ nunit-xml-parser -file random.xml
[
  {
    "Name": "Full.Test.Name",
    "Result": "Success",
    "TimeStarted": "2017-02-15T23:36:16Z",
    "TimeEnded": "2017-02-15T23:36:17Z",
    "DurationSeconds": 1.14,
    "ErrorMessage": null,
    "StackTrace": null
  }
]
```
