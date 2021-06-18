package util_test

import (
	"testing"

	"github.com/cgacurrency/gocga/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCGAAmount(t *testing.T) {
	for _, s := range []string{
		"133246497.546603",
		"1000.000000",
		"0.100000",
		"0.000001",
	} {
		n, err := util.CGAAmountFromString(s)
		require.Nil(t, err)
		assert.Equal(t, s, n.String())
	}
	for _, s := range []string{
		"0.0000000000000000000000000000001",
	} {
		_, err := util.CGAAmountFromString(s)
		require.NotNil(t, err)
	}
}
