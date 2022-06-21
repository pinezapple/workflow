package core

import (
	"io/ioutil"
	"os"
)

func GetFileSize(path string) (filesize int64, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return
	}

	return fi.Size(), err
}

func GetAllFileSizeInDirectory(path string) (filenames []string, filesizes []int64, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, f := range files {
		name := f.Name()
		size, er := GetFileSize(path + "/" + name)
		if err != nil {
			return nil, nil, er
		} else {
			filesizes = append(filesizes, size)
			filenames = append(filenames, path+"/"+name)
		}
	}

	return
}
