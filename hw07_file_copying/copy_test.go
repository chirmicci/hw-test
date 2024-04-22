package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("copy /dev/somepath", func(t *testing.T) {
		err := Copy("/dev/somepath", "/tmp", 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)
	})
	t.Run("file doesn't exist", func(t *testing.T) {
		err := Copy("file.txt", "/tmp", 0, 0)
		require.NotNil(t, err)
	})
	t.Run("copy file to itself", func(t *testing.T) {
		f, _ := os.Create("test.txt")
		defer os.Remove("test.txt")
		_, _ = f.WriteString("test")

		err := Copy("test.txt", "test.txt", 0, 0)
		require.NotNil(t, err)
	})
	t.Run("copy directory", func(t *testing.T) {
		err := Copy("/tmp", "file.txt", 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)
	})
}
