package ip

import (
	"errors"
	"net"
	"strings"
)

type IP struct {
	input string
	ip    net.IP
	net   *net.IPNet
}

func Parse(ip string) (*IP, error) {
	if strings.Contains(ip, "/") {
		address, net, err := net.ParseCIDR(ip)
		if err != nil {
			return nil, err
		}

		return &IP{
			input: ip,
			ip:    address,
			net:   net,
		}, nil
	}

	address := net.ParseIP(ip)
	if address == nil {
		return nil, errors.New("invalid IP address")
	}

	return &IP{
		input: ip,
		ip:    address,
		net:   nil,
	}, nil
}

func (t *IP) String() string {
	return t.input
}
