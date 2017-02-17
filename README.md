# tests-parser

A parser for several test results formats. The tool dumps to JSON.

# usage

`-file {file} -format {format}` will output an array of test results defined in the results file of the specified format.

```
$ tests-parser -file -format nunit random.xml
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

# formats

The following formats are supported:

- NUnit XML
- Jenkins junit plugin
