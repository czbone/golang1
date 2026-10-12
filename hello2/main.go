package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	printHello(os.Stdout)
}

func printHello(w io.Writer) {
	fmt.Fprintln(w, "Hello World")
}
