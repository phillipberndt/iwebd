// UPnP media sharing.
//
// This package relies on Anacrolix's DMS server to serve media via UPnP/DLNA.
package dlnad

import (
	_ "io"
	"log"
	"net"
	"path/filepath"
	"time"

	"github.com/anacrolix/dms/dlna/dms"
	"github.com/spf13/cobra"

	"go.pberndt.com/iwebd/util"
)

// Launch a DLNA server
func Serve(cmd *cobra.Command, args []string) {
	// Generate listener
	addr, _ := cmd.Flags().GetString("addr")
	listener, err := util.GenListener(addr, 80)
	if err != nil {
		util.Log.Error("Failed to serve dlna: %v", err)
		return
	}
	util.Log.Info("%s server starting", "dlna")
	util.ShowInterfaces(listener.Addr())

	// Announce via Bonjour
	if b, _ := cmd.Flags().GetBool("bonjour"); b {
		sq := util.RegisterBonjour("_dlna._upnp._tcp", []string{}, listener.Addr())
		defer sq.Shutdown()
	}

	// Configure DMS
	no_ffmpeg, _ := cmd.Flags().GetBool("no-ffmpeg")
	_, ipNet, _ := net.ParseCIDR("0.0.0.0/0")
	path, _ := filepath.Abs("./")
	s := &dms.Server{
		RootObjectPath: path,
		NotifyInterval: 30 * time.Second,
		AllowedIpNets: []*net.IPNet{ipNet},
		NoTranscode: no_ffmpeg,
		NoProbe: no_ffmpeg,
		HTTPConn: listener,
	}
	log.Default().SetOutput(util.Log)
	log.Default().SetFlags(0)

	// Run DMS
	err = s.Init()
	if err != nil {
		util.Log.Error("Failed to serve dlna: %v", err)
		return
	}

	quitNotifier := make(chan int)
	go func() {
		err := s.Run()
		if err != nil {
			util.Log.Error("Failure while serving: %v", err)
		}
		quitNotifier <- 0;
	}()

	util.WaitToQuit(quitNotifier)

	s.Close()
}
