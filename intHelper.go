package commonUtils

import (
	"fmt"
	"sort"
)

// IntSliceToString - convert a slice to a printable string.
func IntSliceToString(src []int) (tgt string) {
	tgt ="["
	for i:=0; i<len(src);i++ {
		joinChr := ", "
		if i == 0 {
			joinChr = ""
		}
		tgt = fmt.Sprintf("%s%s%d", tgt,joinChr,src[i])
	}
	tgt = fmt.Sprintf("%s]", tgt)
	return
}

// SortDescendingIntSlice - take in a slice and sort it in descending order
func SortDescendingIntSlice(s []int) []int {
	sort.Sort(sort.Reverse(sort.IntSlice(s)))
	return s
}
// SortAscendingIntSlice - take in a slice and sort it in ascending order
func SortAscendingIntSlice(s []int) []int {
	sort.Ints(s)
	return s
}

