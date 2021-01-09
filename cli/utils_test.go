package cli

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScanSplit(t *testing.T) {
	s := `a b "c d e" f`
	arg := ArgSplit(s)
	assert.Equal(t, 4, len(arg))
	assert.Equal(t, []string{"a", "b", "c d e", "f"}, arg)
}
