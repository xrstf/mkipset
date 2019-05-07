package ip

import "sort"

type Set map[string]IP

func NewSet(items ...IP) Set {
	list := make(Set)

	for _, item := range items {
		list[item.String()] = item
	}

	return list
}

func (s Set) Add(ip IP) {
	s[ip.String()] = ip
}

func (s Set) Remove(ip IP) {
	delete(s, ip.String())
}

func (s Set) Sorted() Slice {
	result := make(Slice, len(s))
	idx := 0

	for _, ip := range s {
		result[idx] = ip
		idx++
	}

	sort.Sort(result)

	return result
}
