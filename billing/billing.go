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

package billing

import (
	"bufio"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/imharshita/fair-billing/pkg/stack"
)

const (
	Start = "Start"
	End   = "End"
)

var layout string = "15:04:05"

type sessionReport struct {
	stack         *stack.Stack
	NumOfSessions int
	TotalDuration float64
}

func isFound(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func isValid(str string) bool {
	statusMarkers := []string{Start, End}
	logLine := strings.Split(str, " ")
	if len(logLine) == 3 {
		sessionTimeString := logLine[0]
		_, err := time.Parse(layout, sessionTimeString)
		if err != nil {
			return false
		}
		username := logLine[1]
		isalnum := regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString
		if !isalnum(username) {
			return false
		}
		marker := logLine[2]
		if !isFound(statusMarkers, marker) {
			return false
		}
		return true
	}
	return false
}

func sortedKeys(m map[string]*sessionReport) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// Takes log file as arg and returns sorted keys(for consistency in report), user session report, error (if any)
func Process(file *os.File) ([]string, map[string]*sessionReport, error) {
	var err error
	var earliestTime, latestTime, sessionTime time.Time

	scanner := bufio.NewScanner(file)

	report := make(map[string]*sessionReport)

	lineCounter := 1

	for scanner.Scan() {
		str := strings.Trim(scanner.Text(), " ")
		if isValid(str) {

			logLine := strings.Split(str, " ")
			sessionTimeString := logLine[0]
			sessionTime, err = time.Parse(layout, sessionTimeString)
			if err != nil {
				return nil, nil, err
			}

			username := logLine[1]
			marker := logLine[2]

			if lineCounter == 1 {
				earliestTime = sessionTime
			}

			if _, ok := report[username]; !ok {
				s := stack.New()
				rep := &sessionReport{
					stack: s,
				}
				report[username] = rep
				if marker == Start {
					rep.stack.Push(marker, sessionTime)
				} else if marker == End {
					rep.TotalDuration += sessionTime.Sub(earliestTime).Seconds()
					rep.NumOfSessions += 1
				}
			} else {
				rep := report[username]
				if marker == Start {
					rep.stack.Push(marker, sessionTime)
				} else if marker == End {
					if rep.stack.Len() == 0 {
						rep.TotalDuration += sessionTime.Sub(earliestTime).Seconds()
						rep.NumOfSessions += 1
					} else {
						_, time_stamp := rep.stack.Pop()
						rep.TotalDuration += sessionTime.Sub(time_stamp).Seconds()
						rep.NumOfSessions += 1
					}
				}
			}
		}

		lineCounter++
	}

	latestTime = sessionTime
	for _, value := range report {
		for value.stack.Len() != 0 {
			_, time_stamp := value.stack.Pop()
			value.TotalDuration += latestTime.Sub(time_stamp).Seconds()
			value.NumOfSessions += 1
		}
	}
	// for report consistency
	keys := sortedKeys(report)
	return keys, report, nil
}
