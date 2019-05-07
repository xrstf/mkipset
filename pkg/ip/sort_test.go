package ip

import (
	"testing"
)

func TestSetSorted(t *testing.T) {
	set := NewSet(
		p("1.2.3.4"),
		p("10.2.3.4"),
		p("2.3.4.5"),
	)

	testSort(t, set.Sorted(), Slice{
		p("1.2.3.4"),
		p("2.3.4.5"),
		p("10.2.3.4"),
	})
}
