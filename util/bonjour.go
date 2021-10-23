package util

import (
	"os"
	"net"
	"strconv"
	"strings"

	"github.com/oleksandr/bonjour"
)

// Wrapper around the Bonjour implementation
type BonjourService struct {
	s []*bonjour.Server
}

func (s *BonjourService) Shutdown() {
	for _, sv := range s.s {
		sv.Shutdown()
	}
	Log.Info("Announcing shutdown via Bonjour")
}

// Register a new service to be announced via Bonjour.
//
// Arguments are the Bonjour service name, additional text parameters, and
// the listener of the service to allow this function to decide which
// interfaces to announce the service on.
//
func RegisterBonjour(service string, txt []string, listenerAddr net.Addr) (svr *BonjourService) {
	svr = &BonjourService{}

	host, _ := os.Hostname()
	if i := strings.Index(host, "."); i > 0 {
		host = host[:i]
	}

	shost, sport, _ := net.SplitHostPort(listenerAddr.String())
	port, _ := strconv.Atoi(sport)
	var filterAddr string
	if !net.ParseIP(shost).IsUnspecified() {
		filterAddr = shost
	}

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, v := range addrs {
			s := v.String()
			if strings.Contains(s, ":") {
				continue
			}
			s = s[:strings.Index(s, "/")]
			if filterAddr != "" && strings.Index(s, filterAddr) < 0 {
				continue
			}

			sv, err := bonjour.RegisterProxy("iwebd on " + host, service, "", port, host, s, txt, &i)
			if err != nil {
				Log.Error("Failed to announce via Bonjour on %s: %v", i.Name, err)
				continue
			}
			Log.Info("Announcing via Bonjour on %v with IP %s", i.Name, s)
			svr.s = append(svr.s, sv)
			break
		}
	}

	return
}
