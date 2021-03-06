package util

import (
	"io/ioutil"
	"os"
	"path"
	"pushtart/logging"
	"strings"
)

//GetFilenameListInFolder returns a list of files (not directories) with a given suffix in a given folder. The returned []string
//represents files with an absolute path.
func GetFilenameListInFolder(folder, suffix string) ([]string, error) {
	output := []string{}

	pwd, err := os.Getwd()
	if err != nil {
		logging.Error("file-util", err)
		return nil, err
	}

	files, err := ioutil.ReadDir(path.Join(pwd, folder))
	if err != nil {
		logging.Error("file-util", err)
		return nil, err
	}

	for _, file := range files {
		if (!file.IsDir()) && strings.HasSuffix(file.Name(), suffix) {

			p := file.Name()
			if !path.IsAbs(p) {
				p = path.Join(path.Join(pwd, folder), file.Name())
			}
			output = append(output, p)
		}
	}
	return output, nil
}

//DirExists returns true if a directory exists at the given path.
func DirExists(path string) (bool, error) {
	s, err := os.Stat(path)
	if err == nil {
		if s.IsDir() {
			return true, nil
		}
	}
	return false, err
}

//FileExists returns true if a file exists at the given path.
func FileExists(path string) (bool, error) {
	s, err := os.Stat(path)
	if err == nil {
		if !s.IsDir() {
			return true, nil
		}
	}
	return false, err
}
