package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
}

func SaveFile(b []byte, dName, fName string) (err error) {
	fPath := filepath.Join(dName, fName)
	fPath = filepath.Clean(fPath)
	fs, er := os.Stat(dName)
	if er != nil {
		if os.IsNotExist(er) {
			if err = os.MkdirAll(dName, 0755); err != nil {
				return
			}
		} else {
			return er
		}
	} else {
		if !fs.IsDir() {
			return fmt.Errorf("%s is not a directory", dName)
		}
	}
	obj, err := os.Create(fPath)
	if err != nil {
		return
	}
	defer obj.Close()
	if _, err = obj.Write(b); err != nil {
		return
	}
	return
}
