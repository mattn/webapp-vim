// +build windows

package main

import (
	"bytes"
	"fmt"
	"io"
	"syscall"

	"github.com/mattn/go-encoding"
	enc "golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

var (
	kernel32               = syscall.NewLazyDLL("kernel32")
	procGetConsoleOutputCP = kernel32.NewProc("GetConsoleOutputCP")
	inputEncoding          enc.Encoding
)

func init() {
	r1, _, _ := procGetConsoleOutputCP.Call()
	if r1 != 0 {
		inputEncoding = encoding.GetEncoding(fmt.Sprintf("CP%d", +int(r1)))
	}
}

func convert_input(input []byte) []byte {
	if inputEncoding != nil {
		in := bytes.NewReader(input)
		var out bytes.Buffer
		io.Copy(&out, transform.NewReader(in, inputEncoding.NewDecoder()))
		return out.Bytes()
	}
	return input
}
