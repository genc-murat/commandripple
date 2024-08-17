package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Tree(args []string) error {
	root := "."
	if len(args) > 0 {
		root = args[0]
	}
	return printTree(root, "")
}

func printTree(path string, prefix string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	fmt.Println(prefix + fileInfo.Name())

	if !fileInfo.IsDir() {
		return nil
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for i, file := range files {
		newPrefix := prefix + "├── "
		if i == len(files)-1 {
			newPrefix = prefix + "└── "
		}
		err := printTree(filepath.Join(path, file.Name()), newPrefix)
		if err != nil {
			return err
		}
	}

	return nil
}
