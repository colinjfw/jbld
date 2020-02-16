package queue

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWaits(t *testing.T) {
	c, err := Run(1, []string{"1"}, func(f string) ([]string, error) {
		time.Sleep(1 * time.Millisecond)
		return []string{"2"}, nil
	})
	require.NoError(t, err)
	require.Equal(t, 2, c)
}

func TestLargeSingle(t *testing.T) {
	c, err := Run(1, []string{"1"}, func(f string) ([]string, error) {
		time.Sleep(1 * time.Millisecond)
		return []string{"2", "3", "4", "5", "6"}, nil
	})
	require.NoError(t, err)
	require.Equal(t, 6, c)
}

func TestQueue(t *testing.T) {
	c, err := Run(10, []string{"2", "3", "6"}, func(f string) ([]string, error) {
		time.Sleep(1 * time.Millisecond)
		return []string{"1", "4", "5"}, nil
	})
	require.NoError(t, err)
	require.Equal(t, 6, c)
}

func TestQueueError(t *testing.T) {
	q := New(10, func(f string) ([]string, error) {
		time.Sleep(1 * time.Millisecond)
		return []string{"1", "4", "5"}, errors.New("foobar")
	})
	q.Run([]string{"2", "3", "6"})
}

func TestQueue100Times(t *testing.T) {
	for i := 0; i < 100; i++ {
		c, err := Run(10, []string{"2", "3", "6"}, func(f string) ([]string, error) {
			n := rand.Intn(10)
			time.Sleep(time.Duration(n) * time.Millisecond)
			return []string{"1", "4", "5"}, nil
		})
		require.NoError(t, err)
		require.Equal(t, 6, c)
	}
}
