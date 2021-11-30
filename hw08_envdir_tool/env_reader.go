package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envs := make(map[string]EnvValue)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		f, err := os.Open(fmt.Sprintf("%s/%s", dir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("file open failed: %w", err)
		}

		reader := bufio.NewReader(f)
		row, err := reader.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			err = nil
		}
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("read process failed: %w", err)
		}

		row = bytes.ReplaceAll(row, []byte{0x00}, []byte("\n"))
		fName := strings.ReplaceAll(file.Name(), "=", "")

		str := bytes.TrimRight(row, "\n\t ")

		var envVal EnvValue
		envVal.Value = string(str)
		if len(str) == 0 {
			envVal.NeedRemove = true
		}
		envs[fName] = envVal
	}

	return envs, nil
}
