package ip

import "bytes"

type Slice []IP

// Len implements the sort.Interface
func (s Slice) Len() int {
	return len(s)
}

// Swap implements the sort.Interface
func (s Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implements the sort.Interface
func (s Slice) Less(i, j int) bool {
	a := s[i]
	b := s[j]

	aV4 := a.ip.To4() != nil
	bV4 := b.ip.To4() != nil

	// sort ipv4 always before ipv6
	if aV4 != bV4 {
		return aV4 // i<j if the left side is IPv4
	}

	if a.net == nil || b.net == nil {
		return bytes.Compare(a.ip, b.ip) == -1
	}

	return bytes.Compare(a.net.Mask, b.net.Mask) == -1
}
