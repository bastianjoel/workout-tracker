package appassets

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmbedded(t *testing.T) {
	c, err := FS().Open("dumbbell.png")
	require.NoError(t, err)

	s, err := c.Stat()
	require.NoError(t, err)

	require.NotZero(t, s.Size())
}
