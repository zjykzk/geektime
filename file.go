package geektimedl

import "os"

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func makeSureDirExist(dir string) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil && os.IsExist(err) {
		return nil
	}
	return err
}
