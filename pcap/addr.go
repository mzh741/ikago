package pcap

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type IPVersionOption int

const IPv4AndIPv6 IPVersionOption = 0
const IPv4Only IPVersionOption = 1
const IPv6Only IPVersionOption = 2

// IPPort describes a network endpoint with an IP and a port
type IPPort struct {
	IP              net.IP
	Port            uint16
	IsPortUndefined bool
}

func (i IPPort) String() string {
	var ip string
	if i.IP.To4() != nil {
		ip = i.IP.String()
	} else {
		ip = fmt.Sprintf("[%s]", i.IP)
	}
	if i.IsPortUndefined {
		return ip
	}
	return fmt.Sprintf("%s:%d", ip, i.Port)
}

// ParseIPPort returns an IPPort by the given string of address
func ParseIPPort(s string) (*IPPort, error) {
	if s[0] == '[' {
		// Guess IPv6
		if s[len(s)-1] == ']' {
			// Just IP
			ip := net.ParseIP(s[1 : len(s)-1])
			if ip == nil {
				return nil, fmt.Errorf("parse ip port: %w", fmt.Errorf("invalid ipv6 ip %s", s))
			}
			return &IPPort{
				IP:              ip,
				Port:            0,
				IsPortUndefined: true,
			}, nil
		} else {
			// IP and port
			strs := strings.Split(s[1:], "]:")
			if len(strs) != 2 {
				return nil, fmt.Errorf("parse ip port: %w", fmt.Errorf("invalid ipv6 address %s", s))
			}
			ip := net.ParseIP(strs[0])
			if ip == nil {
				return nil, fmt.Errorf("parse ip port: %w", fmt.Errorf("invalid ipv6 ip %s", strs[0]))
			}
			port, err := strconv.ParseUint(strs[1], 10, 16)
			if err != nil {
				return nil, fmt.Errorf("parse ip port: %w", fmt.Errorf("invalid port %s", strs[1]))
			}
			return &IPPort{
				IP:              ip,
				Port:            uint16(port),
				IsPortUndefined: false,
			}, nil
		}
	} else {
		// Guess IPv4
		if strings.Contains(s, ":") {
			// IP and port
			strs := strings.Split(s, ":")
			if len(strs) != 2 {
				return nil, fmt.Errorf("parse ip port: %w", fmt.Errorf("invalid ipv4 address %s", s))
			}
			ip := net.ParseIP(strs[0])
			if ip == nil {
				return nil, fmt.Errorf("parse ip port: %w", fmt.Errorf("invalid ipv4 ip %s", strs[0]))
			}
			port, err := strconv.ParseUint(strs[1], 10, 16)
			if err != nil {
				return nil, fmt.Errorf("parse ip port: %w", fmt.Errorf("invalid port %s", strs[1]))
			}
			return &IPPort{
				IP:              ip,
				Port:            uint16(port),
				IsPortUndefined: false,
			}, nil
		} else {
			// Just IP
			ip := net.ParseIP(s)
			if ip == nil {
				return nil, fmt.Errorf("parse ip port: %w", fmt.Errorf("invalid ipv4 ip %s", s))
			}
			return &IPPort{
				IP:              ip,
				Port:            0,
				IsPortUndefined: true,
			}, nil
		}
	}
}
