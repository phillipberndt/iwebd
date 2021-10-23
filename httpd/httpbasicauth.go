package httpd

import (
	"net/http"

	"go.pberndt.com/iwebd/util"
)

// Wrap a http handler to require basic authentication.
//
// The realm is hard-coded to iwebd.
func WithAuth(h http.Handler, auth *util.UserPass) (http.Handler) {
	if auth == nil || !auth.IsSet() {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()
		if ok && user == auth.User() && password == auth.Password() {
			h.ServeHTTP(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="iwebd", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
    })
}
