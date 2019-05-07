package blacklist

import (
	"testing"
	"time"

	"github.com/xrstf/mkipset/pkg/ip"
)

func p(str string) *ip.IP {
	parsed, _ := ip.Parse(str)

	return parsed
}

func d(str string) *time.Time {
	parsed, _ := time.Parse("2006-01-02", str)

	return &parsed
}

func TestMergeEntriesWithNoDates(t *testing.T) {
	entries := make(Entries, 0)
	newList := Entries{
		Entry{
			IP:     p("127.0.0.1"),
			After:  nil,
			Before: nil,
		},
	}

	merged := entries.Merge(newList)

	if len(merged) != 1 {
		t.Fatalf("List should contain 1 element, but contains %d.", len(merged))
	}
}

func TestMergeIdenticalEntries(t *testing.T) {
	entries := Entries{
		Entry{
			IP:     p("127.0.0.1"),
			After:  nil,
			Before: nil,
		},
	}

	newList := Entries{
		Entry{
			IP:     p("127.0.0.1"),
			After:  nil,
			Before: nil,
		},
	}

	merged := entries.Merge(newList)

	if len(merged) != 1 {
		t.Fatalf("List should contain 1 element, but contains %d.", len(merged))
	}
}

func TestMergeNonIdenticalEntries(t *testing.T) {
	entries := Entries{
		Entry{
			IP:     p("127.0.0.1"),
			After:  nil,
			Before: nil,
		},
	}

	newList := Entries{
		Entry{
			IP:     p("127.0.0.2"),
			After:  nil,
			Before: nil,
		},
	}

	merged := entries.Merge(newList)

	if len(merged) != 2 {
		t.Fatalf("List should contain 1 element, but contains %d.", len(merged))
	}
}

func TestMergeEntriesWithDates(t *testing.T) {
	entries := Entries{
		Entry{
			IP:     p("127.0.0.1"),
			After:  nil,
			Before: nil,
		},
	}

	newList := Entries{
		Entry{
			IP:     p("127.0.0.1"),
			After:  d("2019-01-01"),
			Before: nil,
		},
	}

	merged := entries.Merge(newList)

	if len(merged) != 1 {
		t.Fatalf("List should contain 1 element, but contains %d.", len(merged))
	}

	if merged[0].After != nil {
		t.Fatal("`after` timestamp should have been ignored and left as nil.")
	}
}

func TestMergeEntriesWithDates2(t *testing.T) {
	entries := Entries{
		Entry{
			IP:     p("127.0.0.1"),
			After:  d("2019-01-01"),
			Before: nil,
		},
	}

	newList := Entries{
		Entry{
			IP:     p("127.0.0.1"),
			After:  nil,
			Before: nil,
		},
	}

	merged := entries.Merge(newList)

	if len(merged) != 1 {
		t.Fatalf("List should contain 1 element, but contains %d.", len(merged))
	}

	if merged[0].After != nil {
		t.Fatal("`after` timestamp should have been reset to nil.")
	}
}

func TestMergeEntriesWithDates3(t *testing.T) {
	target := d("2018-01-01")

	entries := Entries{
		Entry{
			IP:     p("127.0.0.1"),
			After:  d("2019-01-01"),
			Before: nil,
		},
	}

	newList := Entries{
		Entry{
			IP:     p("127.0.0.1"),
			After:  target,
			Before: nil,
		},
	}

	merged := entries.Merge(newList)

	if len(merged) != 1 {
		t.Fatalf("List should contain 1 element, but contains %d.", len(merged))
	}

	if !merged[0].After.Equal(*target) {
		t.Fatalf("`after` timestamp should have been reset to %v.", *target)
	}
}

func TestMergeEntriesWithDates4(t *testing.T) {
	target := d("2018-01-01")

	entries := Entries{
		Entry{
			IP:     p("127.0.0.1"),
			After:  target,
			Before: nil,
		},
	}

	newList := Entries{
		Entry{
			IP:     p("127.0.0.1"),
			After:  d("2019-01-01"),
			Before: nil,
		},
	}

	merged := entries.Merge(newList)

	if len(merged) != 1 {
		t.Fatalf("List should contain 1 element, but contains %d.", len(merged))
	}

	if !merged[0].After.Equal(*target) {
		t.Fatalf("`after` timestamp should have been kept at %v.", *target)
	}
}
