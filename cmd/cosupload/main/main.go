package main

import (
	"fmt"
	"os"
	"tccos-tools/cmd/cosupload"
	"tccos-tools/exitcode"
)

func main() {
	if e := cosupload.Cmd().Execute(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
	exitcode.Exit()
}
