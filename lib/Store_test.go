package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	store := NewStore()

	// Test empty store
	{
		messages := store.FindAll()
		require.Equal(t, len(messages), 0)
	}

	// Test populated store
	{
		expected := "Hello, world!"
		store.Save(expected)
		messages := store.FindAll()
		require.Equal(t, 1, len(messages))
		require.Equal(t, expected, messages[0])
	}

	// TODO: Test for race condition?
}
