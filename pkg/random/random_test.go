package random_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/pkg/random"
)

func TestRandomString_Success(t *testing.T) {

	s := random.String(10, true, true, true)
	log.Printf(s)
	assert.Equal(t, 10, len(s))
}

func TestRandomString_Different(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	exists := map[string]bool{}

	for i := 0; i < 1000000; i++ {
		s := random.String(10, true, true, true)
		_, ok := exists[s]
		assert.False(t, ok)
		exists[s] = true
	}
}

func TestRandomString_Different_OrderID(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	exists := map[string]bool{}

	for i := 0; i < 1000000; i++ {
		s := random.String(8, false, true, true)
		_, ok := exists[s]
		assert.False(t, ok)
		exists[s] = true
	}
}
