package logfile

import (
	"os"
	"path"
)

func Create(filepath string) (*os.File, error) {
	// create the necessary folders in case they
	// do not already exist
	folder := path.Dir(filepath)
	os.MkdirAll(folder, os.ModePerm)

	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}
