package commonUtils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapStringIntToString(t *testing.T) {
	src := map[string]int{
		"test": 1,
		"all":  2,
		"things": 3,
	}
	actual := MapStringIntToString(src)

	expected := "[\"test\": 1, \"all\": 2, \"things\": 3]"
	assert.Equal(t, expected, actual)
}
