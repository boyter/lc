package main

import "path/filepath"
import "os"

const dirPath = "/home/bboyter/"

func main() {
	// walk all files in directory
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			println(info.Name())
		}
		return nil
	})
}
