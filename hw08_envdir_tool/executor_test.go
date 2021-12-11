package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	args := []string{"testdata/echo.sh", "key1=1", "key2=2", "key3=3"}
	envs, _ := ReadDir("testdata/env")

	t.Run("valid data", func(t *testing.T) {
		n := RunCmd(args, envs)

		require.Equal(t, 0, n)
	})

	t.Run("invalid data", func(t *testing.T) {
		n := RunCmd([]string{"tratatata/echo.sh"}, envs)

		require.Equal(t, 1, n)
	})
}
