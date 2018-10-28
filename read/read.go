package read

import (
	"errors"
	"io/ioutil"
	"path"
)

// SM returns a byte array from an SM file
func SM(smPath string) ([]uint8, error) {
	if path.Ext(smPath) == ".sm" {
		data, _ := ioutil.ReadFile(smPath)
		return data, nil
	}
	return nil, errors.New("Extension Error: File is not of type .sm")
}
