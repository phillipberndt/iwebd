// A good old FTP server.
//
// This package relies on the goftp.io FTP server.
//
package ftpd

import (
	"fmt"
	"strconv"
	"strings"

	"goftp.io/server/v2"
	"goftp.io/server/v2/driver/file"
	"github.com/spf13/cobra"

	"go.pberndt.com/iwebd/util"
)

// Authentication scheme that accepts each user/password combination
type zeroAuth struct{}
func (i *zeroAuth) CheckPasswd(ctx *server.Context, u string, p string) (bool, error) {
	return true, nil
}

// Launch a FTP server
func Serve(cmd *cobra.Command, args []string) {
	// Generate Listener
	addr, _ := cmd.Flags().GetString("addr")
	listener, err := util.GenListener(addr, 21)
	if err != nil {
		util.Log.Error("Failed to serve ftp: %v", err)
		return
	}
	listenerAddr := listener.Addr().String()
	listenerParts := strings.Split(listenerAddr, ":")
	port, err := strconv.Atoi(listenerParts[len(listenerParts) - 1])
	if err != nil {
		panic("Unexpected problem parsing port")
	}

	util.Log.Info("%s server starting", "ftp")
	util.ShowInterfaces(listener.Addr())

	// Define FTP server
	auth, _ := cmd.Flags().Lookup("auth").Value.(*util.UserPass);
	readOnly, _ := cmd.Flags().GetBool("read-only")
	var driver server.Driver
	if readOnly {
		driver, err = newReadOnlyDriver("./")
	} else {
		driver, err = file.NewDriver("./")
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	var serverAuth server.Auth
	if auth.IsSet() {
		serverAuth = &server.SimpleAuth{Name: auth.User(), Password: auth.Password()}
	} else {
		serverAuth = &zeroAuth{}
	}

	opt := &server.Options{
		Name:   "iwebd",
		Driver: driver,
		Port:   port,
		Auth:   serverAuth,
		Perm:   newFsPerm("./"),
		Logger: newLogger(),
		WelcomeMessage: "Welcome to iwebd powered by goftp.io",
	}

	s, err := server.NewServer(opt)
	if err != nil {
		util.Log.Error("Failed to serve ftp: %v", err)
		return
	}

	// Announce via Bonjour
	if b, _ := cmd.Flags().GetBool("bonjour"); b {
		sq := util.RegisterBonjour("_ftp._tcp", []string{}, listener.Addr())
		defer sq.Shutdown()
	}

	// Run server
	quitNotifier := make(chan int)
	go func() {
		err = s.Serve(listener)
		if err != nil && err != server.ErrServerClosed {
			util.Log.Error("Failed to serve ftp: %v", err)
		}
		quitNotifier <- 0;
	}()

	util.WaitToQuit(quitNotifier)

	s.Shutdown()
}
