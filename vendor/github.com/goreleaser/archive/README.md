# archive

A simple Go archiving library.

This project adheres to the Contributor Covenant [code of conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.
We appreciate your contribution. Please refer to our [contributing guidelines](CONTRIBUTING.md).

[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Travis](https://img.shields.io/travis/goreleaser/archive.svg?style=flat-square)](https://travis-ci.org/goreleaser/archive)
[![Coverage Status](https://img.shields.io/codecov/c/github/goreleaser/archive/master.svg?style=flat-square)](https://codecov.io/gh/goreleaser/archive)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/goreleaser/archive)
[![Go Report Card](https://goreportcard.com/badge/github.com/goreleaser/archive?style=flat-square)](https://goreportcard.com/report/github.com/goreleaser/archive)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)

## Example usage

```go
file, err := os.Create("file.zip")
if err != nil {
  // deal with the error
}
archive := archive.New(file)
defer archive.Close()
archive.Add("file.txt", "/path/to/file.txt")
```
