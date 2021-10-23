package util

import (
	"net"
	"os"
	"strings"
	"strconv"
)

// Output information about how to reach a listener.
func ShowInterfaces(listenerAddr net.Addr) {
	name, _ := os.Hostname()

	shost, sport, _ := net.SplitHostPort(listenerAddr.String())
	port, _ := strconv.Atoi(sport)
	var filterAddr string
	if !net.ParseIP(shost).IsUnspecified() {
		filterAddr = shost
	}

	Log.Info("iwbed running on %s, serving on port %d for", name, port)

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		var filteredAddrs []string
		for _, v := range addrs {
			s := v.String()
			if strings.Contains(s, ":") {
				continue
			}
			s = s[:strings.Index(s, "/")]
			if filterAddr != "" && strings.Index(s, filterAddr) < 0 {
				continue
			}
			filteredAddrs = append(filteredAddrs, s);
		}
		if len(filteredAddrs) == 0 {
			continue
		}
		Log.Info(" %s: %s", i.Name, strings.Join(filteredAddrs, ", "))
	}
}
