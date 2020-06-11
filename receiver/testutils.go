package main

import (
	"io"
	"io/ioutil"
	"strings"
)

// readCloser converts a string to an io.ReadCloser using
// the iotuil.NopCloser. This is a testing utility
func readCloser(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}
