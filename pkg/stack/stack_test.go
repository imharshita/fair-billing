package stack

import (
	"testing"
	"time"
)

var layout string = "15:04:05"

func Test(t *testing.T) {
	s := New()

	time_stamp, _ := time.Parse(layout, layout)

	if s.Len() != 0 {
		t.Errorf("Length of an empty stack should be 0")
	}

	s.Push("one", time_stamp)

	if s.Len() != 1 {
		t.Errorf("Length should be 0")
	}

	str, tm := s.Pop()
	if str != "one" || tm != time_stamp {
		t.Errorf("Top items on the stack should be %v, %v", str, tm)
	}

	if s.Len() != 0 {
		t.Errorf("Stack should be empty")
	}

	s.Push("one", time_stamp)
	s.Push("two", time_stamp)

	if s.Len() != 2 {
		t.Errorf("Length should be 2")
	}
}
