package ipset

import (
	"sort"

	"github.com/xrstf/mkipset/pkg/ip"
)

type Interface interface {
	Create(setname string, t SetType) error
	Add(setname string, entry string) error
	Delete(setname string, entry string) error
	List(setname string) error
	Swap(oldname string, newname string) error
	Destroy(setname string) error
}

type Set struct {
	Name     string
	Type     SetType
	Revision int
	Header   SetHeader
	Members  []string
}

func (s Set) MembersEquals(others ip.Slice) bool {
	ours := make(ip.Slice, 0)

	for _, member := range s.Members {
		parsed, err := ip.Parse(member)
		if err == nil {
			ours = append(ours, *parsed)
		}
	}

	if len(others) != len(ours) {
		return false
	}

	sort.Sort(ours)
	sort.Sort(others)

	for idx, our := range ours {
		if !our.Equals(others[idx]) {
			return false
		}
	}

	return true
}

type SetHeader struct {
	Family      string
	HashSize    int
	MaxElements int
	MemorySize  int
	References  int
}
