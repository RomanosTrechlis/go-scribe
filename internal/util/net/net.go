package net

import (
	"net"
	"os"
	"strings"
)

// GetIPAddress finds the first ip address after the loopback and returns it.
// In case there is an error returns the os.Hostname().
// In the future should become more sophisticated, but it's ok for now.
func GetIPAddress() (string, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return os.Hostname()
	}
	address := addr[len(addr)-1].String()
	if len(addr) > 3 {
		address = addr[len(addr)-3].String()
	}
	addressArr := strings.Split(address, "/")
	if len(addressArr) > 2 {
		return os.Hostname()
	}

	return addressArr[0], nil
}
