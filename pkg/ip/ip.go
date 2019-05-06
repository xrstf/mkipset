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

func (i *IP) String() string {
	return i.input
}

func (i *IP) IsNet() bool {
	return i.net != nil
}

func (i *IP) Collides(other *IP) bool {
	if i.IsNet() {
		if other.IsNet() {
			return i.net.Contains(other.net.IP) || other.net.Contains(i.net.IP)
		} else {
			return i.net.Contains(other.ip)
		}
	} else {
		if other.IsNet() {
			return other.net.Contains(i.ip)
		} else {
			return i.ip.Equal(other.ip)
		}
	}
}
