package ports

import (
	"strconv"
	"strings"
)

const (
	RaftServerGRPC = 14250
)

// PortToHostPort converts the port into a host:port address string
func PortToHostPort(port int) string {
	return ":" + strconv.Itoa(port)
}

// FormatHostPort returns hostPort in a usable format (host:port) if it wasn't already
func FormatHostPort(hostPort string) string {
	if hostPort == "" {
		return ""
	}

	if strings.Contains(hostPort, ":") {
		return hostPort
	}

	return ":" + hostPort
}
