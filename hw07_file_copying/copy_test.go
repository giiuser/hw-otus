package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	for _, tcase := range [...]struct {
		in        string
		out       string
		benchmark string
		offset    int64
		limit     int64
	}{
		{
			in:        "./testdata/input.txt",
			out:       "out_offset0_limit0.txt",
			benchmark: "./testdata/out_offset0_limit0.txt",
			offset:    0,
			limit:     0,
		},
		{
			in:        "./testdata/input.txt",
			out:       "out_offset0_limit10.txt",
			benchmark: "./testdata/out_offset0_limit10.txt",
			offset:    0,
			limit:     10,
		},
		{
			in:        "./testdata/input.txt",
			out:       "out_offset0_limit1000.txt",
			benchmark: "./testdata/out_offset0_limit1000.txt",
			offset:    0,
			limit:     1000,
		},
		{
			in:        "./testdata/input.txt",
			out:       "out_offset0_limit10000.txt",
			benchmark: "./testdata/out_offset0_limit10000.txt",
			offset:    0,
			limit:     10000,
		},
		{
			in:        "./testdata/input.txt",
			out:       "out_offset100_limit1000.txt",
			benchmark: "./testdata/out_offset100_limit1000.txt",
			offset:    100,
			limit:     1000,
		},
		{
			in:        "./testdata/input.txt",
			out:       "out_offset6000_limit1000.txt",
			benchmark: "./testdata/out_offset6000_limit1000.txt",
			offset:    6000,
			limit:     1000,
		},
	} {
		t.Run(fmt.Sprintf("test-%q", tcase.out), func(t *testing.T) {
			_ = Copy(tcase.in, tcase.out, tcase.offset, tcase.limit)
			out, _ := ioutil.ReadFile(tcase.out)
			benchmark, _ := ioutil.ReadFile(tcase.benchmark)
			defer os.Remove(tcase.out)

			if !bytes.Equal(out, benchmark) {
				t.Errorf("incoming file and outcomming file not matched")
			}
		})
	}

	t.Run("File unknown length", func(t *testing.T) {
		err := Copy("/dev/urandom", "output.txt", 0, 0)
		expected := errors.New("unsupported file")
		require.Equal(t, expected, err)
	})

	t.Run("Offset more file size", func(t *testing.T) {
		err := Copy("./testdata/input.txt", "output.txt", 10000, 50)
		expected := errors.New("offset exceeds file size")
		require.Equal(t, expected, err)
	})
}
