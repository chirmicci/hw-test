package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("unset env", func(t *testing.T) {
		os.Setenv("TEST", "tttt")
		env := make(Environment)
		env["TEST"] = EnvValue{
			"clever",
			false,
		}
		RunCmd([]string{"ls"}, env)
		test, ok := os.LookupEnv("TEST")
		if !ok {
			fmt.Println("env not found")
		}
		require.Equal(t, "", test)
	})

	t.Run("empty cmd & env", func(t *testing.T) {
		var s []string
		r := RunCmd(s, Environment{})
		require.Equal(t, -1, r)
	})

	t.Run("empty env", func(t *testing.T) {
		r := RunCmd([]string{"ls"}, Environment{})
		require.Equal(t, 0, r)
	})
}
