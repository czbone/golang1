package main

import (
	"bytes"
	"testing"
)

func TestPrintHello(t *testing.T) {
	var output bytes.Buffer

	printHello(&output)

	if got, want := output.String(), "Hello World\n"; got != want {
		t.Errorf("printHello() = %q, want %q", got, want)
	}
}
