package commonUtils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortDescendingIntSlice(t *testing.T) {
	actual := []int{0,7,8,6,2,5}
	expecting := []int{8,7,6,5,2,0}
	actual = SortDescendingIntSlice(actual)
	fmt.Printf("Compare: %s to \n         %s\n",
		IntSliceToString(expecting), IntSliceToString(actual))
	assert.Equal(t,expecting, actual)
}

func TestSortAscendingIntSlice(t *testing.T) {
	tmp := []int{0,7,8,6,2,5}
	expecting := []int{0,2,5,6,7,8}
	actual := SortAscendingIntSlice(tmp)
	fmt.Printf("Compare: %s to \n         %s\n",
		IntSliceToString(expecting), IntSliceToString(actual))
	assert.Equal(t,expecting, actual)
}