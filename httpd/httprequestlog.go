package httpd

import (
	"net"
	"net/http"

	"go.pberndt.com/iwebd/util"
)

type statusCatcher struct {
    http.ResponseWriter
    Status int
}

func (r *statusCatcher) WriteHeader(status int) {
    r.Status = status
    r.ResponseWriter.WriteHeader(status)
}

// Wrap a HTTP handler to emit request logs using
// the logger in go.pberndt.com/iwebd/util.Log.
func WithHttpLogging(h http.Handler) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        s := &statusCatcher{
            ResponseWriter: w,
        }
        h.ServeHTTP(s, r)

		remoteHost, _, _ := net.SplitHostPort(r.RemoteAddr)
		util.Log.Info("%s %d %s", remoteHost, s.Status, r.URL.Path)
    })
}
