/*
Copyright 2021.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
		if pErr, ok := err.(*os.PathError); ok {
			fmt.Fprintln(out, "Failed to open file at path", pErr.Path)
			os.Exit(1)
		}
		fmt.Fprintln(out, "Generic error", err)
		os.Exit(1)
	}
	defer file.Close()
	keys, report, err := billing.Process(file)
	if err != nil {
		fmt.Fprintln(out, err)
		os.Exit(1)
	} else {
		for _, k := range keys {
			value := report[k]
			fmt.Fprintln(out, k, value.NumOfSessions, value.TotalDuration)
		}
	}
}
