package main

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"os"
	"path"
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
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment, len(files))
	if len(files) == 0 {
		return env, nil
	}
	var fileOpen *os.File
	defer fileOpen.Close()

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), "=") {
			err = errors.New("file name contains '='")
			return nil, err
		}
		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		if info.Size() == 0 {
			env[file.Name()] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}

		fileOpen, err := os.Open(path.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		reader := bufio.NewReader(fileOpen)
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Print(err)
			return nil, err
		}

		value := string(bytes.ReplaceAll(line, []byte{0x00}, []byte("\n")))
		value = strings.TrimRight(value, " \t")

		env[file.Name()] = EnvValue{
			Value: value,
		}
	}

	return env, nil
}
