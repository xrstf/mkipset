package blacklist

import (
	"sort"
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

func (e Entries) RemoveCollisions(ips ip.Slice) Entries {
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

func (e Entries) Merge(others Entries) Entries {
	result := e

	for _, other := range others {
		var existing *Entry

		for idx, self := range e {
			if self.IP.Equals(*other.IP) {
				existing = &e[idx]
				break
			}
		}

		if existing == nil {
			result = append(result, other)
			continue
		}

		// compare time stamps
		// we only ever extend any given time range, assuming that a
		// disapparing block is the worse outcome of a merge

		existingAfter := existing.After
		otherAfter := other.After

		if otherAfter == nil || (existingAfter != nil && otherAfter.Before(*existingAfter)) {
			existing.After = otherAfter
		}

		existingBefore := existing.Before
		otherBefore := other.Before

		if otherBefore == nil || (existingBefore != nil && otherBefore.After(*existingBefore)) {
			existing.Before = otherBefore
		}
	}

	return result
}

func (e Entries) IPs() ip.Slice {
	result := make(ip.Slice, 0)

	for _, entry := range e {
		if entry.IP != nil {
			result = append(result, *entry.IP)
		}
	}

	sort.Sort(result)

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
