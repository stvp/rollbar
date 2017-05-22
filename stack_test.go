package rollbar

import (
	"testing"
)

func TestBuildStack(t *testing.T) {
	frame := BuildStack(1)[0]
	if frame.Filename != "github.com/stvp/rollbar/stack_test.go" {
		t.Errorf("got: %s", frame.Filename)
	}
	if frame.Method != "rollbar.TestBuildStack" {
		t.Errorf("got: %s", frame.Method)
	}
	if frame.Line != 8 {
		t.Errorf("got: %d", frame.Line)
	}
}

func TestShortenFilePath(t *testing.T) {
	tests := []struct {
		Given    string
		Expected string
	}{
		{"", ""},
		{"foo.go", "foo.go"},
		{"/usr/local/go/src/pkg/runtime/proc.c", "pkg/runtime/proc.c"},
		{"/home/foo/go/src/github.com/stvp/rollbar.go", "github.com/stvp/rollbar.go"},
	}
	for i, test := range tests {
		got := shortenFilePath(test.Given)
		if got != test.Expected {
			t.Errorf("tests[%d]: got %s", i, got)
		}
	}
}
