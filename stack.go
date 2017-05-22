package rollbar

import (
	"os"
	"runtime"
	"strings"
)

var (
	knownFilePathPatterns = []string{
		"github.com/",
		"code.google.com/",
		"bitbucket.org/",
		"launchpad.net/",
	}
)

// Frame is a single line of executed code in a Stack.
type Frame struct {
	Filename string `json:"filename"`
	Method   string `json:"method"`
	Line     int    `json:"lineno"`
}

// Stack represents a stacktrace as a slice of Frames.
type Stack []Frame

// BuildStack builds a full stacktrace for the current execution location.
func BuildStack(skip int) Stack {
	stack := make(Stack, 0)

	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		file = shortenFilePath(file)
		stack = append(stack, Frame{file, functionName(pc), line})
	}

	return stack
}

// Remove un-needed information from the source file path. This makes them
// shorter in Rollbar UI as well as making them the same, regardless of the
// machine the code was compiled on.
//
// Examples:
//   /usr/local/go/src/pkg/runtime/proc.c -> pkg/runtime/proc.c
//   /home/foo/go/src/github.com/rollbar/rollbar.go -> github.com/rollbar/rollbar.go
func shortenFilePath(s string) string {
	idx := strings.Index(s, "/src/pkg/")
	if idx != -1 {
		return s[idx+5:]
	}
	for _, pattern := range knownFilePathPatterns {
		idx = strings.Index(s, pattern)
		if idx != -1 {
			return s[idx:]
		}
	}
	return s
}

func functionName(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "???"
	}
	name := fn.Name()
	end := strings.LastIndex(name, string(os.PathSeparator))
	return name[end+1 : len(name)]
}
