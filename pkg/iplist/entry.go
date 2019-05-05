package iplist

import "time"

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

type Entry struct {
	IP     string
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
