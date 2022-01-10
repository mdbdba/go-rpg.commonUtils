package commonUtils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringSliceToString(t *testing.T) {
	src := []string{
		"test",
		"all",
		"things",
	}
	actual := StringSliceToString(src)

	expected := "[\"test\", \"all\", \"things\"]"
	assert.Equal(t, expected, actual)
}