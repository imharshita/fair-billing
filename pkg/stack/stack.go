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

// This package is inspired by https://pkg.go.dev/github.com/go-stack/stack
package stack

import "time"

type (
	Stack struct {
		top    *node
		length int
	}
	node struct {
		value      string
		prev       *node
		time_stamp time.Time
	}
)

// Create a new stack
func New() *Stack {
	return &Stack{nil, 0}
}

// Return the number of items in the stack
func (this *Stack) Len() int {
	return this.length
}

// Pop the top items of the stack and return it
func (this *Stack) Pop() (string, time.Time) {
	if this.length == 0 {
		var t time.Time
		return "", t
	}

	n := this.top
	this.top = n.prev
	this.length--
	return n.value, n.time_stamp
}

// Push items onto the top of the stack
func (this *Stack) Push(value string, time_stamp time.Time) {
	n := &node{value, this.top, time_stamp}
	this.top = n
	this.length++
}
