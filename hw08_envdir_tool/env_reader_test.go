package main

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

type _file struct {
	key, value string
}

func createMockDir() (string, error) {
	name, err := os.MkdirTemp("./", "mock")
	if err != nil {
		return "", err
	}
	return name, nil
}

func createFiles(d string, f []_file) error {
	for _, v := range f {
		file, err := os.Create(path.Join(d, v.key))
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(v.value); err != nil {
			return err
		}
	}
	return nil
}

func TestReadDir(t *testing.T) {
	t.Run("create & read files", func(t *testing.T) {
		dest, err := createMockDir()
		if err != nil {
			fmt.Println(err)
		}
		files := []_file{
			{"A", "simple"},
			{"B", "width new line \n blabla"},
			{"D", "width \x00terminal\000nulls"},
			{"E", ""},
		}

		if err := createFiles(dest, files); err != nil {
			fmt.Println(err)
		}

		result, err := ReadDir(dest)

		require.Nil(t, err)
		require.Equal(t, Environment{
			"A": EnvValue{"simple", false},
			"B": EnvValue{"width new line", false},
			"D": EnvValue{"width \nterminal\nnulls", false},
			"E": EnvValue{"", true},
		}, result)

		if err := os.RemoveAll(dest); err != nil {
			fmt.Println("can`t clear tmp dir")
		}
	})
	t.Run("wrong file name", func(t *testing.T) {
		dest, err := createMockDir()
		if err != nil {
			fmt.Println(err)
		}
		files := []_file{
			{"B=", "width new line \n blabla"},
		}

		if err := createFiles(dest, files); err != nil {
			fmt.Println(err)
		}

		_, err = ReadDir(dest)

		require.NotNil(t, err)

		if err := os.RemoveAll(dest); err != nil {
			fmt.Println("can`t clear tmp dir")
		}
	})
}
