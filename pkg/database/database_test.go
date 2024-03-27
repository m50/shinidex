package database

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateId(t *testing.T) {
	id := generateId()
	assert.Equal(t, 32, len(id))
}

func TestGenerateIdSequential(t *testing.T) {
	ids := make([]string, 5)
	for i := 0; i < 5; i++ {
		ids[i] = generateId()
	}
	assert.True(t, sort.StringsAreSorted(ids), "Strings are not sorted")
}
