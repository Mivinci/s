package model

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

type Foo struct {
	A string
	B int
}

func TestIsZeroOrInvalid(t *testing.T) {
	s1 := Foo{
		A: "a",
	}
	assert(t, IsZeroOrInvalid(s1, "A") == nil, "")
	assert(t, IsZeroOrInvalid(s1, "A", "B") != nil, "")
	assert(t, IsZeroOrInvalid(s1, "A", "C") != nil, "")
}

func TestIsZeroOrInvalidPtr(t *testing.T) {
	s1 := &Foo{
		A: "a",
	}
	assert(t, IsZeroOrInvalid(s1, "A") == nil, "")
	assert(t, IsZeroOrInvalid(s1, "A", "B") != nil, "")
	assert(t, IsZeroOrInvalid(s1, "A", "C") != nil, "")
}
