package ip

import (
	"sort"
	"testing"
)

func p(ip string) IP {
	parsed, _ := Parse(ip)

	return *parsed
}

func testSort(t *testing.T, input Slice, expected Slice) {
	sort.Sort(input)

	for i, result := range input {
		expectation := expected[i]

		if !result.Equals(expectation) {
			t.Fatal("not sorted")
		}
	}
}

func TestIPv4BeforeIPv6(t *testing.T) {
	testSort(t, Slice{
		p("::1"),
		p("127.0.0.1"),
	}, Slice{
		p("127.0.0.1"),
		p("::1"),
	})
}

func TestNumericSort(t *testing.T) {
	testSort(t, Slice{
		p("1.2.3.4"),
		p("10.2.3.4"),
		p("2.3.4.5"),
	}, Slice{
		p("1.2.3.4"),
		p("2.3.4.5"),
		p("10.2.3.4"),
	})
}

func TestNumericSortMasks(t *testing.T) {
	testSort(t, Slice{
		p("1.2.3.4/1"),
		p("1.2.3.4/10"),
		p("1.2.3.4/2"),
		p("1.2.3.4/32"),
		p("1.2.3.4/4"),
	}, Slice{
		p("1.2.3.4/1"),
		p("1.2.3.4/2"),
		p("1.2.3.4/4"),
		p("1.2.3.4/10"),
		p("1.2.3.4/32"),
	})
}
