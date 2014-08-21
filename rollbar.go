package rollbar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/adler32"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	NAME    = "go-rollbar"
	VERSION = "0.1.0"

	// Severity levels
	CRIT  = "critical"
	ERR   = "error"
	WARN  = "warning"
	INFO  = "info"
	DEBUG = "debug"
)

var (
	// Rollbar access token. If this is blank, no errors will be reported to
	// Rollbar.
	Token = ""

	// All errors and messages will be submitted under this environment.
	Environment = "development"

	// API endpoint for Rollbar.
	Endpoint = "https://api.rollbar.com/api/1/item/"

	// Maximum number of errors allowed in the sending queue before we start
	// dropping new errors on the floor.
	Buffer = 1000

	// Queue of messages to be sent.
	bodyChannel chan map[string]interface{}
	waitGroup   sync.WaitGroup
)

// -- Setup

func init() {
	bodyChannel = make(chan map[string]interface{}, Buffer)

	go func() {
		for body := range bodyChannel {
			post(body)
			waitGroup.Done()
		}
	}()
}

// -- Error reporting

// Error asynchronously sends an error to Rollbar with the given severity level.
func Error(level string, err error) {
	ErrorWithStackSkip(level, err, 1)
}

// ErrorWithStackSkip asynchronously sends an error to Rollbar with the given
// severity level and a given number of stack trace frames skipped.
func ErrorWithStackSkip(level string, err error, skip int) {
	parts := strings.SplitN(err.Error(), "\n", 2)
	body := buildBody(level, parts[0])
	data := body["data"].(map[string]interface{})
	errBody, fingerprint := errorBody(err, skip)
	data["body"] = errBody
	data["fingerprint"] = fingerprint

	push(body)
}

// -- Message reporting

// Message asynchronously sends a message to Rollbar with the given severity
// level. Rollbar request is asynchronous.
func Message(level string, msg string) {
	parts := strings.SplitN(msg, "\n", 2)
	body := buildBody(level, parts[0])
	data := body["data"].(map[string]interface{})
	data["body"] = messageBody(msg)

	push(body)
}

// -- Misc.

// Wait will block until the queue of errors / messages is empty.
func Wait() {
	waitGroup.Wait()
}

// Build the main JSON structure that will be sent to Rollbar with the
// appropriate metadata.
func buildBody(level, title string) map[string]interface{} {
	timestamp := time.Now().Unix()
	hostname, _ := os.Hostname()

	return map[string]interface{}{
		"access_token": Token,
		"data": map[string]interface{}{
			"environment": Environment,
			"title":       title,
			"level":       level,
			"timestamp":   timestamp,
			"platform":    runtime.GOOS,
			"language":    "go",
			"server": map[string]interface{}{
				"host": hostname,
			},
			"notifier": map[string]interface{}{
				"name":    NAME,
				"version": VERSION,
			},
		},
	}
}

// Build an error inner-body for the given error. If skip is provided, that
// number of stack trace frames will be skipped.
func errorBody(err error, skip int) (map[string]interface{}, string) {
	stack := BuildStack(3 + skip)
	fingerprint := stack.Fingerprint()
	errBody := map[string]interface{}{
		"trace": map[string]interface{}{
			"frames": stack,
			"exception": map[string]interface{}{
				"class":   errorClass(err),
				"message": err.Error(),
			},
		},
	}
	return errBody, fingerprint
}

// Build a message inner-body for the given message string.
func messageBody(s string) map[string]interface{} {
	return map[string]interface{}{
		"message": map[string]interface{}{
			"body": s,
		},
	}
}

func errorClass(err error) string {
	class := reflect.TypeOf(err).String()
	if class == "" {
		return "panic"
	} else if class == "*errors.errorString" {
		checksum := adler32.Checksum([]byte(err.Error()))
		return fmt.Sprintf("{%x}", checksum)
	} else {
		return strings.TrimPrefix(class, "*")
	}
}

// -- POST handling

// Queue the given JSON body to be POSTed to Rollbar.
func push(body map[string]interface{}) {
	if len(bodyChannel) < Buffer {
		waitGroup.Add(1)
		bodyChannel <- body
	} else {
		stderr("buffer full, dropping error on the floor")
	}
}

// POST the given JSON body to Rollbar synchronously.
func post(body map[string]interface{}) {
	if len(Token) == 0 {
		stderr("empty token")
		return
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		stderr("failed to encode payload: %s", err.Error())
		return
	}

	resp, err := http.Post(Endpoint, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		stderr("POST failed: %s", err.Error())
	} else if resp.StatusCode != 200 {
		stderr("received response: %s", resp.Status)
	}
	if resp != nil {
		resp.Body.Close()
	}
}

// -- stderr

func stderr(format string, args ...interface{}) {
	format = "Rollbar error: " + format + "\n"
	fmt.Fprintf(os.Stderr, format, args...)
}
