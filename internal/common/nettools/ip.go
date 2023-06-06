// Package nettools provide utility functions to work with network.
package nettools

import (
	"errors"
	"net"
)

const RealIPHeader = "X-Real-IP"

var ErrHostIPFailed = errors.New("unable to define host IP")

// GetHostIP - retrieves host IP.
func GetHostIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			// interface down
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			// loopback interface
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				// not an ipv4 address
				continue
			}
			return ip, nil
		}
	}

	return nil, ErrHostIPFailed
}
