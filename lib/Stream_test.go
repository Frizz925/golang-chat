package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {
	stream := NewStream()
	ch, id := stream.Subscribe()
	expected := "Hello, world"

	// Test stream firing
	{
		timeout := time.After(500 * time.Millisecond)
		go stream.Fire(expected)

		select {
		case <-timeout:
			t.Fatal("Test timeout while waiting for event firing")
		case message := <-ch:
			require.Equal(t, expected, message)
		}
	}

	// Test stream shutdown
	{
		timeout := time.After(500 * time.Millisecond)
		go stream.Shutdown()

		select {
		case <-timeout:
			t.Fatal("Test timeout while waiting for shutdown")
		case message := <-ch:
			require.Empty(t, message)
		}
	}

	// Test stream unsubscribe
	{
		timeout := time.After(500 * time.Millisecond)
		stream.Unsubscribe(id)
		go stream.Fire(expected)

		select {
		case <-ch:
			t.Fatal("Expected no event fired after unsubscribing")
		case <-timeout:
		}
	}
}
