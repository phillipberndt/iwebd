// Provides a simple logging interface for console applications.
//
// No log level configuration, log line support, file output, etc. is supported.
// But the logs automatically apply coloring to formatted elements in Printf() like
// statements, which is easy to use and looks nice.
package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
)

// A Logger is an object for logging to the console. An application would usually want
// to instanciate this once, using NewLogger(), and then use that instance throughout
// its application to log.
//
// Logger implements Writer to allow attaching it to other logging modules' outputs.
type Logger struct {
	hasColor bool
}

// Create a new Logger.
func NewLogger() (*Logger) {
	return &Logger{isatty.IsTerminal(os.Stdout.Fd())}
}

func (l *Logger) rotateColor(i int) string {
	colors := []string{"\033[36m", "\033[32m", "\033[33m", "\033[34m"} // etc
	if l.hasColor {
		return colors[i % len(colors)]
	} else {
		return ""
	}
}

func (l *Logger) defaultColor() string {
	if l.hasColor {
		return "\033[39m"
	} else {
		return ""
	}
}

func (l *Logger) errorColor() string {
	if l.hasColor {
		return "\033[31m"
	} else {
		return ""
	}
}

func (l *Logger) boldFormatting() string {
	if l.hasColor {
		return "\033[1m"
	} else {
		return ""
	}
}

func (l *Logger) resetFormatting() string {
	if l.hasColor {
		return "\033[m"
	} else {
		return ""
	}
}

// Main entry point for the Logger.
//
// This is the function that does the heavy lifting.
// It outputs lines as follows to the standard output:
//
// prefix " [" timestamp "] " mid fmt
//
// where fmt is the output of Sprint(f, a...), with color codes injected
// around the verbs.
//
func (l *Logger) Log(prefix string, f string, mid string, a ...interface{}) {
	now := time.Now().Format(time.RFC3339)

	args := make([]interface{}, len(a))
	fstr := strings.Builder{}

	// Timestamp
	fstr.WriteString(prefix)
	fstr.WriteString("[")
	fstr.WriteString(l.rotateColor(0))
	fstr.WriteString(now)
	fstr.WriteString(l.defaultColor())
	fstr.WriteString("] ")
	fstr.WriteString(mid)

	// Each argument
	colCtr := 0
	for i := strings.Index(f, "%"); i >= 0; i = strings.Index(f, "%") {
		if i == len(f) - 1 || f[i+1] == '%' {
			f = f[i+2:]
			fstr.WriteString("%%")
			continue
		}

		fstr.WriteString(f[:i])
		colCtr += 1
		fstr.WriteString(l.rotateColor(colCtr))
		fstr.WriteString("%")
		f = f[i+1:]
		n := strings.IndexAny(f, "FvExOtbqoescgdpfGTUX")
		if n < 0 {
			break
		}
		fstr.WriteString(f[:n+1])
		f = f[n+1:]
		fstr.WriteString(l.defaultColor())
	}
	fstr.WriteString(f)
	for i, v := range a {
		args[i] = v
	}

	fstr.WriteString(l.resetFormatting())
	fstr.WriteString("\n")

	fmt.Printf(fstr.String(), args[:len(a)]...)
}

// Logs with Info severity.
func (l *Logger) Info(f string, a ...interface{}) {
	l.Log("", f, "", a...)
}

// Logs with Error severity.
//
// Output still goes to stdout, but there is a red "[ERROR]" in the output.
func (l *Logger) Error(f string, a ...interface{}) {
	l.Log(l.boldFormatting(), f, l.errorColor() + "[ERROR] " + l.defaultColor(), a...)
}

// Implements Writer
func (l *Logger) Write(p []byte) (n int, err error) {
	ps := string(p)
	ps = strings.TrimSpace(ps)
	l.Log("", ps, "")
	return len(p), nil
}
