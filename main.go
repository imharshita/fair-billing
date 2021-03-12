package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/imharshita/fair-billing/billing"
)

var out io.Writer = os.Stdout

func main() {
	numOfArgs := len(os.Args)
	if numOfArgs < 2 {
		fmt.Fprintln(out, errors.New("Log file path not provided"))
		os.Exit(1)
	} else if numOfArgs > 2 {
		fmt.Fprintln(out, errors.New("Too may arguments"))
		os.Exit(1)
	}
	fileName := os.Args[1]
	fileName = filepath.FromSlash(fileName)
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintln(out, err)
	}
	defer file.Close()
	keys, report, err := billing.Process(file)
	if err != nil {
		fmt.Fprintln(out, err)
	} else {
		for _, k := range keys {
			value := report[k]
			fmt.Fprintln(out, k, value.NumOfSessions, value.TotalDuration)
		}
	}
}
