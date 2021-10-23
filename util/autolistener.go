package util

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func GenListener(address string, defaultPort uint16) (net.Listener, error) {
	parts := strings.Split(address, ":")
	var (
		host string
		port uint16
		portStr string
	)
	if len(parts) == 1 && parts[0] != "" {
		portStr = parts[0]
	} else if len(parts) == 2 {
		host = parts[0]
		portStr = parts[1]
	} else {
		portStr = "0"
	}
	u64port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, err
	}
	port = uint16(u64port)

	if port != 0 {
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
		return listener, err
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, defaultPort))
	if err == nil {
		return listener, nil
	}

	for port = 1234; port < 10000; port++ {
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
		if err == nil {
			return listener, nil
		}
	}

	return nil, err
}
