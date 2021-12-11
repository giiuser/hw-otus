package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	testDir := "testdata/env"

	t.Run("valid data", func(t *testing.T) {
		if err := os.WriteFile("testdata/env/CUSTOM=", []byte("custom"), 0o644); err != nil {
			log.Fatal(err)
		}
		envs, err := ReadDir(testDir)
		require.NoError(t, err)
		require.Equal(t, "bar", envs["BAR"].Value)
		require.Equal(t, "", envs["EMPTY"].Value)
		require.Equal(t, "custom", envs["CUSTOM"].Value)
		if err := os.Remove("testdata/env/CUSTOM="); err != nil {
			log.Fatal(err)
		}
	})
}
