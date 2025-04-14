package crontask

import (
	"errors"
	"testing"
)

func TestErrAllTypes(t *testing.T) {
	// Call the newErr method with various types
	e := newErr(
		"stringTest",
		[]string{"array", "of", "strings"},
		rune(':'), // Just joins without additional space
		42,
		3.14,
		true,
		errors.New("customError"),
	)

	expected := "stringTest array of strings: 42 3.14 true customError"

	if e.Error() != expected {
		t.Errorf("got: %q, expected: %q", e.Error(), expected)
	}
}
