package dirutil

import (
	"fmt"
	"os"
)

func MdAllIfNeeded(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		err := os.MkdirAll(path, os.ModePerm)
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("file exists")
	}
	return nil
}
