package iplist

import (
	"time"

	"github.com/xrstf/mkipset/pkg/ip"
)

type Entries []Entry

func (e Entries) Active(now time.Time) Entries {
	result := make(Entries, 0)

	for _, entry := range e {
		if entry.IsActive(now) {
			result = append(result, entry)
		}
	}

	return result
}

func (e Entries) RemoveCollisions(ips []ip.IP) Entries {
	result := make(Entries, 0)

	for _, entry := range e {
		conflicts := false

		for _, i := range ips {
			if i.Collides(entry.IP) {
				conflicts = true
				break
			}
		}

		if !conflicts {
			result = append(result, entry)
		}
	}

	return result
}

func (e Entries) IPs() []ip.IP {
	result := make([]ip.IP, 0)

	for _, entry := range e {
		if entry.IP != nil {
			result = append(result, *entry.IP)
		}
	}

	return result
}

type Entry struct {
	IP     *ip.IP
	After  *time.Time
	Before *time.Time
}

func (e *Entry) IsActive(now time.Time) bool {
	if e.After != nil && now.Before(*e.After) {
		return false
	}

	if e.Before != nil && now.After(*e.Before) {
		return false
	}

	return true
}
