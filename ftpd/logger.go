package ftpd

import (
	"go.pberndt.com/iwebd/util"
)

type cmdPair struct {
	command string
	params string
}

// Logger for ftp server
type logger struct{
	last map[string]cmdPair
}

// Returns a Logger for goftp that forwards logs
// to iwebd's internal logging system.
func newLogger() *logger {
	var rv logger
	rv.last = make(map[string]cmdPair)
	return &rv
}

func (l *logger) Print(sessionID string, message interface{}) {
	util.Log.Info("%s %v", sessionID, message)
}

func (l *logger) Printf(sessionID string, format string, v ...interface{}) {
	args := append([]interface{}{sessionID}, v...)

	util.Log.Info("%s " + format, args...)
}
func (l *logger) PrintCommand(sessionID string, command string, params string) {
	l.last[sessionID] = cmdPair{command, params}
}
func (l *logger) PrintResponse(sessionID string, code int, message string) {
	cmd := l.last[sessionID]
	if code / 100 > 3 /* Error */ || code / 10 == 25 /* File modification */ {
		util.Log.Info("%s %-5s %-20s → %d %s", sessionID, cmd.command, cmd.params, code, message)
	}
}
