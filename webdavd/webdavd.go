// webdav file access.
//
package webdavd

import (
	"context"
	"net/http"
	"crypto/tls"
	"time"

	"golang.org/x/net/webdav"
	"github.com/spf13/cobra"

	"go.pberndt.com/iwebd/httpd"
	"go.pberndt.com/iwebd/util"
)

// Launch a webdav(s) server.
func Serve(cmd *cobra.Command, args []string) {
	useTls, _ := cmd.Flags().GetBool("tls")
	addr, _ := cmd.Flags().GetString("addr")
	auth, _ := cmd.Flags().Lookup("auth").Value.(*util.UserPass);
	readOnly, _ := cmd.Flags().GetBool("read-only")

	var fs webdav.FileSystem
	if readOnly {
		fs = readOnlyDir{"./"}
	} else {
		fs = webdav.Dir("./")
	}
	h := &webdav.Handler{
			Prefix:     "/",
			FileSystem: fs,
			LockSystem: webdav.NewMemLS(),
			Logger:     nil,
		}

	util.Log.Info("%s server starting", "webdav")

	var defaultPort uint16
	if useTls {
		defaultPort = 443
	} else {
		defaultPort = 80
	}
	listener, err := util.GenListener(addr, defaultPort)
	if err != nil {
		util.Log.Error("Failed to spawn server: %v", err)
		return
	}

	util.ShowInterfaces(listener.Addr())

	var tlsConfig *tls.Config
	if useTls {
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{util.GenKeyPair()},
		}
	}

	s := &http.Server{
		Handler: httpd.WithHttpLogging(httpd.WithAuth(h, auth)),
		TLSConfig: tlsConfig,
	}

	if b, _ := cmd.Flags().GetBool("bonjour"); b {
		sq := util.RegisterBonjour("_webdav._tcp", []string{}, listener.Addr())
		defer sq.Shutdown()
	}

	quitNotifier := make(chan int)
	go func() {
		var err error;
		if !useTls {
			err = s.Serve(listener)
		} else {
			err = s.ServeTLS(listener, "", "")
		}
		if err != nil && err != http.ErrServerClosed {
			util.Log.Error("Failure while serving: %v", err)
		}
		quitNotifier <- 0;
	}()

	util.WaitToQuit(quitNotifier)

	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
	s.Shutdown(ctx)
}
