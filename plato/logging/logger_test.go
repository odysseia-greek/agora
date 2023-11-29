package logging

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

const message = "This is a message"

func TestSystem(t *testing.T) {
	testLogFunction(t, System, message, purple)
}

func TestInfo(t *testing.T) {
	testLogFunction(t, Info, message, green)
}

func TestDebug(t *testing.T) {
	testLogFunction(t, Debug, message, blue)
}

func TestWarn(t *testing.T) {
	testLogFunction(t, Warn, message, yellow)
}

func TestError(t *testing.T) {
	testLogFunction(t, Error, message, red)
}

func TestTrace(t *testing.T) {
	testLogFunction(t, Trace, message, cyan)
}

func testLogFunction(t *testing.T, logFunc func(string), expected string, color string) {
	// Redirect the log output to a buffer
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// Call the log function
	logFunc(expected)

	// Reset log output to the default
	log.SetOutput(nil)

	// Check the output
	output := buf.String()

	assert.Contains(t, output, expected)
	assert.Contains(t, output, color)
}
