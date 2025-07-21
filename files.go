package main

import (
	"os"
	"path/filepath"
	"strings"
)

func getFiles(paths []string) ([]string, error) {
	var fileList []string

	for _, imgpath := range paths {

		file, err := os.Open(imgpath)
		if err != nil {
			return fileList, err
		}

		fi, err := file.Stat()
		file.Close()
		if err != nil {
			return fileList, err
		}

		if fi.IsDir() {
			filepath.Walk(imgpath, func(path string, fInfo os.FileInfo, _ error) error {
				if fInfo.Mode().IsRegular() && !strings.HasPrefix(filepath.Base(path), ".") {
					p, _ := filepath.Abs(path)
					fileList = append(fileList, p)
				}
				return nil
			})

			continue
		}

		if !fi.Mode().IsRegular() {
			continue
		}

		p, _ := filepath.Abs(imgpath)
		fileList = append(fileList, p)
	}
	return fileList, nil
}

func filterFiles(paths []string, exts []string) []string {
	n := 0
pathLoop:
	for _, path := range paths {
		fExt := strings.ToLower(filepath.Ext(path))
		for _, ext := range exts {
			if fExt == ext {
				paths[n] = path
				n++
				continue pathLoop
			}
		}
	}
	return paths[:n]
}
