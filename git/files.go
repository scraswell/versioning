package git

import (
	"errors"
	"os"
)

func directoryExists(filePath string) bool {
	stat, err := os.Stat(filePath)

	if err == nil && !stat.IsDir() {
		return false
	} else if err == nil {
		return true
	} else if err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		panic("stat call failed")
	}
}
