package util

import (
	"os"
	"os/signal"
)

// Helper function to allow servers to wait for either the server quitting or
// the user to interrupt the process, and then still cleaning up.
func WaitToQuit(appQuit <-chan int) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
		case <-c:
			Log.Info("Shutting down")
			return
		case <-appQuit:
			return
	}
}
